// Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
 * @Author: wenchao
 * @Date: 2021-12-25 10:10:42
 * @LastEditors: jeffery
 * @LastEditTime: 2022-03-01 16:00:07
 * @Description:
 */
package docker

import (
	"agent/biz/alivechecker"
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/model/disk_initial/model"
	"agent/config"
	"agent/utils/docker/dockerfacade"
	"agent/utils/hardware"
	"agent/utils/simpleeventbus"
	"agent/utils/tools"
	"time"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

const (
	EventPowerOn = "EventPowerOn"
	EventPairing = "EventPairing"
	EventReset   = "EventReset"

	ContainersUnStarted      = -1
	ContainersWaitOSReady    = -2
	ContainersStarting       = 0
	ContainersStarted        = 1
	ContainersStartedFail    = 2
	ContainersDownloading    = 3
	ContainersDownloaded     = 4
	ContainersDownloadedFail = 5
)

var StartingProgress int // 启动进度, 0-100
var Status int
var ContainerPreUp int
var docker *dockerfacade.DockerFacade

var eventLoopNormal *simpleeventbus.EventLoop
var tickerStartProgress *time.Ticker

func init() {
	writeDefaultDockerComposeFile()
	if hardware.RunningInDocker() {
		writeUpgradeComposeFile()
	}
}

func Start() {
	logger.AppLogger().Debugf("BOXSTATUS===>>>> GetAdminPairedInfo=%+v", clientinfo.GetAdminPairedInfo())

	if !device_ability.GetAbilityModel().InnerDiskSupport &&
		!device_ability.GetAbilityModel().SupportUSBDisk {
		err := MigrateFileStorageData()
		if err != nil {
			logger.AppLogger().Errorf("failed MigrateFileStorageData ,err:%v", err)
			return
		}
	}

	go func() {
		logger.AppLogger().Debugf("Docker.go Entry Start...")
		abilityModel := device_ability.GetAbilityModel()

		disposeComposeFile()

		setDockerPreUp(ContainersUnStarted)
		setDockerStatus(ContainersUnStarted)
		docker = dockerfacade.NewDockerFacade()
		docker.SetClientVersion(config.Config.Docker.APIVersion)
		eventLoopNormal = simpleeventbus.NewEventLoop()

		if !abilityModel.RunInDocker {
			// 磁盘未初始化成功时恢复 docker 相关数据.
			diskInitialInfo := model.ReadDiskInitialInfo()
			logger.AppLogger().Debugf("restoreDockerMetaDir, ReadDiskInitialInfo diskInitialInfo:%+v", diskInitialInfo)
			if device_ability.GetAbilityModel().InnerDiskSupport && // 设备支持外接磁盘，且未初始化成功的状态时
				diskInitialInfo.DiskInitialCode != model.DiskInitialCode_Nomal {
				logger.AppLogger().Debugf("restoreDockerMetaDir, diskInitialInfo.DiskInitialCode != model.DiskInitialCode_Nomal")

				if err := restoreDockerMetaDir(); err != nil {
					logger.AppLogger().Errorf("restoreDockerMetaDir, err:%v", err)
				}
				if err := restoreDockerStorageFile(); err != nil {
					logger.AppLogger().Errorf("restoreDockerStorageFile err:%v", err)
				}
			}

			// 启动 dockerd
			if err := StartDockerEngine(); err != nil {
				logger.AppLogger().Errorf("StartDockerEngine err:%v", err)
				tools.RunCmd("systemctl", []string{"start", "docker.service"})
			}

			waitEngineReady()
		}

		createDockerNetwork()

		go registerHandlerNormal()
		eventLoopNormal.PostEvent(EventPowerOn)
	}()
}

func registerHandlerNormal() {
	logger.AppLogger().Debugf("registerHandlerNormal")
	eventLoopNormal.RegisterEvent(EventPowerOn, func(event string) {
		logger.AppLogger().Debugf("EventPowerOn")
		if clientinfo.HasPairedBefore() {
			logger.AppLogger().Debugf("clientinfo.HasPairedBefore==true")
			ProcessEnv(config.Config.Docker.ComposeFile, nil)
			disposeComposeFile()
			if err := ContainersUpAndPrune(config.Config.Docker.ComposeFile, nil); err != nil {
				logger.AppLogger().Warnf("dockerUpAndPrune err:%v", err)
			} else {
				if GetDockerStatus() == ContainersStarted {
					if !device.GetConfig().EnableInternetAccess {
						StopContainerImmediately(config.Config.Docker.NetworkClientContainerName)
					}
				}
			}
		} else { // 没有配对就没有必要的环境变量参数，没法启动 gateway 等.
			logger.AppLogger().Debugf("clientinfo.HasPairedBefore==false")
			ProcessEnv(config.Config.Docker.ComposeFile, nil)
			disposeComposeFile()
			if device_ability.GetAbilityModel().RunInDocker {
				dockerPull()
			} else {
				dockerPreUp()
			}
		}
		logger.AppLogger().Debugf("EventPowerOn finish")
		PublishDockerPowerOn()
	})

	eventLoopNormal.RegisterEvent(EventPairing, func(event string) {
		logger.AppLogger().Debugf("EventPairing")

		disposeComposeFile()
		ProcessEnv(config.Config.Docker.ComposeFile, nil)
		logger.AppLogger().Debugf("EventPairing, ContainerPreUp=%v", ContainerPreUp)
		if ContainerPreUp == ContainersStarting || ContainerPreUp == ContainersDownloading { // wait...
			waitSeconds := 180
			if device_ability.GetAbilityModel().RunInDocker {
				waitSeconds = 360
			}
			for i := 0; i < waitSeconds; i++ {
				logger.AppLogger().Debugf("EventPairing, waiting loop, ContainerPreUp=%v", ContainerPreUp)
				time.Sleep(time.Second)
				if ContainerPreUp != ContainersStarting && ContainerPreUp != ContainersStartedFail &&
					Status != ContainersDownloading && Status != ContainersDownloaded {
					break
				}
			}
		}
		logger.AppLogger().Debugf("EventPairing, before dockerUpAndPruneWithNoRecreate, ContainerPreUp=%v", ContainerPreUp)
		dockerUpAndPruneWithNoRecreate([]string{config.Config.Docker.NetworkClientContainerName})
		if !device.GetConfig().EnableInternetAccess {
			StopContainerImmediately(config.Config.Docker.NetworkClientContainerName)
		}
		logger.AppLogger().Debugf("EventPairing return")
	})

	eventLoopNormal.RegisterEvent(EventReset, func(event string) {
		logger.AppLogger().Debugf("EventReset")
		dockerDown()
	})
	eventLoopNormal.Poll()
}

func restoreDockerMetaDir() error {
	logger.AppLogger().Debugf("#### restoreDockerStorageFile")

	// 如果初始化磁盘过程中断电，则有可能是 /home/eulixspace_link -> /mnt/bp/...
	// 需要修改成 /home/eulixspace_link -> /home/eulixspace
	// ln -snf /home/eulixspace /home/eulixspace_link
	if _, _, err := tools.RunCmd("ln", []string{"-snf",
		config.Config.Docker.VolumeDirReal,
		config.Config.Docker.VolumeDirLink}); err != nil {
		logger.AppLogger().Warnf("ln -snf %v %v, err:%v",
			config.Config.Docker.VolumeDirReal,
			config.Config.Docker.VolumeDirLink,
			err)
		return err
	}
	return nil
}

// 恢复  /etc/sysconfig/docker-storage
func restoreDockerStorageFile() error {
	logger.AppLogger().Debugf("#### restoreDockerStorageFile")

	file := config.Config.Docker.DockerStorageFile
	backupFile := file + ".backup"

	if fileutil.IsFileExist(backupFile) { // 上次已经备份完成
		logger.AppLogger().Debugf("file exist %v", backupFile)

		if _, _, err := tools.RunCmd("/bin/cp", []string{"-f", backupFile, file}); err != nil {
			logger.AppLogger().Warnf("/bin/cp -f %v %v failed, err:%v", backupFile, file, err)
			return err
		}
	} else {
		s := `# This file may be automatically generated by an installation program.

# By default, Docker uses a loopback-mounted sparse file in
# /var/lib/docker.  The loopback makes it slower, and there are some
# restrictive defaults, such as 100GB max storage.

# If your installation did not set a custom storage for Docker, you
# may do it below.

# Example: Use a custom pair of raw logical volumes (one for metadata,
# one for data).
# DOCKER_STORAGE_OPTIONS = --storage-opt dm.metadatadev=/dev/mylogvol/my-docker-metadata --storage-opt dm.datadev=/dev/mylogvol/my-docker-data

DOCKER_STORAGE_OPTIONS=--graph /userdata/docker --storage-driver overlay2
`

		err := fileutil.WriteToFile(file, []byte(s), true)
		if err != nil {
			logger.AppLogger().Warnf("WriteToFile %v failed, err:%v", file, err)
			return err
		}

	}

	if _, _, err := tools.RunCmd("systemctl", []string{"daemon-reload"}); err != nil {
		logger.AppLogger().Warnf("systemctl stop docker, err:%v", err)
		return err
	}

	return nil
}

func DockerDownImmediately() error {
	return dockerDown()
}

func StopDockerEngine() error {
	if _, _, err := tools.RunCmd("systemctl", []string{"stop", "docker"}); err != nil {
		logger.AppLogger().Warnf("systemctl stop docker, err:%v", err)
		return err
	}
	if _, _, err := tools.RunCmd("/bin/docker", []string{"ps"}); err != nil {
		logger.AppLogger().Debugf("docker ps err:%v", err)
		return err
	}

	logger.AppLogger().Debugf("SUCC docker.StopDockerEngine")
	return nil
}

func StartDockerEngine() error {
	if _, _, err := tools.RunCmd("systemctl", []string{"start", "docker"}); err != nil {
		logger.AppLogger().Warnf("systemctl start docker, err:%v", err)
		return err
	}
	if _, _, err := tools.RunCmd("/bin/docker", []string{"ps"}); err != nil {
		logger.AppLogger().Debugf("docker ps err:%v", err)
		return err
	}

	logger.AppLogger().Debugf("SUCC docker.StartDockerEngine")
	return nil
}

func DockerUpImmediately(tmpEnv map[string]map[string]string) error {
	var err error
	for i := 0; i < int(config.Config.Docker.DockerUpRetryTimes); i++ { // 失败时重试若干次.

		// 未绑定时和绑定时 docker-compose.yml 是不一样的，多了磁盘挂载部分。这样在磁盘初始化、磁盘扩容等场景之后需要根据实际磁盘情况再次处理一下 docker-compose.yml。
		disposeComposeFile()
		ProcessEnv(config.Config.Docker.ComposeFile, tmpEnv)

		err = ContainersUpAndPrune(config.Config.Docker.ComposeFile, nil)
		if err != nil {
			logger.AppLogger().Debugf("#### DockerUpImmediately failed, waiting retry(%v/%v) docker-compose up ...  err:%v",
				i+1, config.Config.Docker.DockerUpRetryTimes, err)
			if i < int(config.Config.Docker.DockerUpRetryTimes) {
				time.Sleep(time.Second * time.Duration(config.Config.Docker.DockerUpRetryIntervalSec))
			}
		} else {
			logger.AppLogger().Debugf("#### DockerUpImmediately succ, stop retry(%v/%v)",
				i+1, config.Config.Docker.DockerUpRetryTimes)
			break
		}
	}
	return err
}

func StopContainerImmediately(containerName string) error {
	return dockerStop(containerName)
}

func setDockerPreUp(preUp int) {
	logger.AppLogger().Debugf("setDockerPreUp, preUp:%+v", preUp)
	ContainerPreUp = preUp
}

func setDockerStatus(started int) {
	logger.AppLogger().Debugf("setDockerStatus, started:%+v", started)
	Status = started
	PublishContainerStaus(Status)
}

func startProgress() {
	StartingProgress = 0
	logger.AppLogger().Debugf("startProgress")
	tickerStartProgress = time.NewTicker(800 * time.Millisecond)
	go timerStartProgress(tickerStartProgress)
}

func timerStartProgress(ticker *time.Ticker) {
	for range ticker.C {
		StartingProgress += 1
		if StartingProgress >= 99 {
			StartingProgress = 99
		}
	}
}

func stopProgress() {
	logger.AppLogger().Debugf("stopProgress")

	if tickerStartProgress != nil {
		tickerStartProgress.Stop()
		tickerStartProgress = nil
	}
}

func finishProgress() {
	logger.AppLogger().Debugf("finishProgress")

	// 时间关系, 以下临时这么处理.
	// compose up -d 已经执行完成, 循环检测网关接口一定次数直到成功.
	tryTotal := 80
	for i := 0; i < tryTotal; i++ {
		if alivechecker.GetContainerStatus(alivechecker.ContainerNameGateway()) {
			logger.AppLogger().Debugf("finishProgress, GetContainerStatus(%v/%v) of %v alive",
				i+1, tryTotal, alivechecker.ContainerNameGateway())
			stopProgress()
			setDockerStatus(ContainersStarted)
			StartingProgress = 100
			break
		} else {
			logger.AppLogger().Debugf("finishProgress, GetContainerStatus(%v/%v) of %v not alive",
				i+1, tryTotal, alivechecker.ContainerNameGateway())
			if i == tryTotal-1 {
				stopProgress()
				setDockerStatus(ContainersStartedFail)
				break
			}

			time.Sleep(time.Second * 3)
		}
	}

}

func GetStartingProgress() int {
	return StartingProgress
}

func GetDockerStatus() int {
	logger.AppLogger().Debugf("GetDockerStatus, DockerStatus:%+v", Status)
	return Status
}

func PostEvent(event string) {
	logger.AppLogger().Debugf("PostEvent, event:%+v", event)
	eventLoopNormal.PostEvent(event)
}
