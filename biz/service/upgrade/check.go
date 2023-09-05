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

package upgrade

import (
	"agent/biz/model/device_ability"
	"agent/biz/model/upgrade"
	"agent/biz/notification"
	"agent/config"
	"agent/utils/docker/dockerfacade"
	"agent/utils/hardware"
	"agent/utils/logger"
	"agent/utils/tools"
	"github.com/dungeonsnd/gocom/file/fileutil"
	"path"
	"strings"
	"time"
)

type UpgradeRecord struct {
	AgentVersion  string `json:"agent_version"`
	RebootPushed  bool   `json:"reboot_pushed"`
	SuccessPushed bool   `json:"success_pushed"`
	Rebooted      bool   `json:"reboot"`
	UpgradeTime   int64  `json:"upgrade_time"`
}

type VersionDownInfo struct {
	VersionId  string    `json:"versionId"`
	Downloaded bool      `json:"downloaded"`
	PkgPath    string    `json:"pkgPath"`
	UpdateTime time.Time `json:"updateTime"`
}

//type Task struct {
//	VersionId        string          `json:"versionId"`
//	Status           string          `json:"status"`        // 整体流程状态："", downloading, downloaded, installing, installed, download-err，install-err
//	DownStatus       string          `json:"downStatus"`    // 下载状态："", ing, done, err
//	InstallStatus    string          `json:"installStatus"` // 安装状态："", ing, done, err
//	StartDownTime    time.Time       `json:"startDownTime"`
//	StartInstallTime time.Time       `json:"startInstallTime"`
//	DoneDownTime     time.Time       `json:"doneDownTime"`
//	DoneInstallTime  time.Time       `json:"doneInstallTime"`
//	RpmPkg           VersionDownInfo `json:"rpmPkg"`
//	CFile            VersionDownInfo `json:"cFile"` // docker-compose.yml
//	ContainerImg     VersionDownInfo `json:"containerImg"`
//}

//const (
//	Ing         = "ing"
//	Done        = "done"
//	Err         = "err"
//	Downloading = "downloading"
//	Downloaded  = "downloaded"
//	Installing  = "installing"
//	Installed   = "installed"
//	DownloadErr = "download-err"
//	InstallErr  = "install-err"
//)

var t upgrade.Task

// 文件监控不知道为什么有情况下监控不到变化. 猜测可能是 upgrade 里面修改 task.json 用了移动文件的方式导致的.
// 所以使用循环读取的方式。

func StartUpgradeMonitor() {
	logger.NotificationLogger().Infof("StartUpgradeMonitor")

	var Dir = path.Join(config.Config.RunTime.BasePath, config.Config.RunTime.DBDir)
	collPath := path.Join(Dir, config.Config.RunTime.UpgradeCollection)
	taskPath := path.Join(collPath, config.Config.RunTime.TaskResource+".json")

	task, err1 := readTask(taskPath)
	if err1 != nil {
		logger.NotificationLogger().Warnf("Failed to readTask, taskPath:%v, err:%+v", taskPath, err1)
	} else {
		logger.NotificationLogger().Infof("task:%+v", task)
	}

	logger.NotificationLogger().Infof("StartUpgradeMonitor, wait now, task:%+v", task)

	for {
		time.Sleep(time.Duration(config.Config.Box.UpgradeCheckIntervalMs) * time.Millisecond)

		if fileutil.IsFileNotExist(taskPath) {
			// logger.NotificationLogger().Debugf("IsFileNotExist taskPath:%+v", taskPath)
			continue
		}

		newTask, err1 := readTask(taskPath)
		if err1 != nil {
			logger.AppLogger().Warnf("Failed to readTask, taskPath:%v, err:%+v", taskPath, err1)
		} else {
			// logger.NotificationLogger().Infof("task:%+v", task)
			// logger.NotificationLogger().Infof("newTask:%+v", newTask)

			if task.Status == upgrade.Downloading && newTask.Status == upgrade.Downloaded { // downloading -> downloaded
				err := notification.OnUpgradeDownloadedSuccess(newTask.VersionId)
				if err != nil {
					logger.NotificationLogger().Warnf("failed onUpgradeDownloadedSuccess, err:%v", err)
				}
			} else if task.Status == upgrade.Downloaded && newTask.Status == upgrade.Installing { // downloaded -> installing
				_, err := notification.OnUpgradeInstalling()
				if err != nil {
					logger.NotificationLogger().Warnf("failed onUpgradeInstalling, err:%v", err)
				}
			}

			task = newTask
		}
	}
}

func readTask(taskPath string) (task upgrade.Task, err error) {
	//task := Task{Status: Installed}
	err = fileutil.ReadFileJsonToObject(taskPath, &task)
	t = task
	return task, err
}

func CheckUpgradeSucc() {
	logger.UpgradeLogger().Infof("current version number from build info: %v", config.VersionNumber)
	current_agent_version := config.VersionNumber
	logger.NotificationLogger().Debugf("CheckUpgradeSucc, current_agent_version:%v", current_agent_version)
	record := &UpgradeRecord{
		AgentVersion:  current_agent_version,
		SuccessPushed: false,
		RebootPushed:  false,
		Rebooted:      false}

	f := config.Config.Notification.UpgradeRecordFile
	if fileutil.IsFileNotExist(f) {
		err := fileutil.WriteToFileAsJson(f, record, "  ", true)
		if err != nil {
			logger.NotificationLogger().Warnf("failed WriteToFileAsJson, err:%v", err)
		} else {
			logger.NotificationLogger().Debugf("succ WriteToFileAsJson, record:%+v", record)
		}
	} else {
		err := fileutil.ReadFileJsonToObject(f, record)
		if err != nil {
			logger.NotificationLogger().Warnf("failed ReadFileJsonToObject, err:%v", err)
		} else {
			logger.NotificationLogger().Debugf("succ ReadFileJsonToObject, obj:%+v", record)
		}
	}

	_, stdout, _ := tools.RunCmd("uptime", []string{"-s"})

	if record.AgentVersion != current_agent_version {
		record.Rebooted = false
		//record.RebootPushed = false
	}

	uptime, err := time.ParseInLocation("2006-01-02 15:04:05", strings.Trim(stdout, "\n"), time.Local)
	logger.UpgradeLogger().Debugf("box uptime:%s", stdout)
	logger.UpgradeLogger().Debugf("time now unix:%d", time.Now().Unix())
	logger.UpgradeLogger().Debugf("up time unix:%d", uptime.Unix())
	if time.Now().Unix()-uptime.Unix() < 60*5 {
		record.Rebooted = true
	}

	// 不需要重启，而且没有推送，而且当前版本和 record 版本不一致
	// 说明不需要重启的升级刚完成，需要推送升级成功
	if !t.NeedReboot && !record.SuccessPushed && record.AgentVersion != current_agent_version {
		logger.NotificationLogger().Infof("do not need to reboot,push success msg to redis")
		err := notification.OnUpgradeSuccess()
		if err != nil {
			logger.NotificationLogger().Warnf("failed onUpgradeSuccess, err:%v", err)
		}
		record.AgentVersion = current_agent_version
		record.SuccessPushed = true
		record.UpgradeTime = time.Now().Unix()
	}
	// 需要重启，且已经重启，而且当前版本和 record 版本不一致
	// 说明需要重启的升级已重启完成，需要推送升级成功，更新record状态
	if t.NeedReboot && record.Rebooted && record.RebootPushed && record.AgentVersion != current_agent_version {
		logger.NotificationLogger().Infof("have been reboot,push success msg to redis")
		err := notification.OnUpgradeSuccess()
		if err != nil {
			logger.NotificationLogger().Warnf("failed onUpgradeSuccess, err:%v", err)
		}
		record.AgentVersion = current_agent_version
		record.SuccessPushed = true
		record.UpgradeTime = time.Now().Unix()
		//record.Rebooted = true
	}
	// 需要重启，还没有重启，而且当前版本和 record 版本不一致
	// 说明需要重启的升级还没有推送重启，重启前需要推送重启消息
	logger.NotificationLogger().Debugf("record.AgentVersion:%s", record.AgentVersion)
	logger.NotificationLogger().Debugf("current_agent_version:%s", current_agent_version)
	logger.NotificationLogger().Debugf("t.NeedReboot:%v", t.NeedReboot)
	logger.NotificationLogger().Debugf("t.NeedReboot && record.AgentVersion != current_agent_version:%v", t.NeedReboot && record.AgentVersion != current_agent_version)
	if t.NeedReboot && record.AgentVersion != current_agent_version {
		logger.NotificationLogger().Infof("need to reboot,push reboot msg to redis")
		if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_SUPPORTED_FROM_MODEL_NUMBER {
			err := notification.OnUpgradeRestart()
			if err != nil {
				logger.NotificationLogger().Warnf("failed onUpgradeRetart, err:%v", err)
			}
			//record.AgentVersion = current_agent_version
			//record.Rebooted = true
			// 发送消息到aosapce-upgrade ，让upgrade重启盒子
			//err = http.SendMessageToSocket(config.Config.RunTime.BasePath+config.Config.RunTime.SocketFile, "reboot")
			//if err != nil {
			//	logger.NotificationLogger().Warnf("failed to send unix socket message, err:%v", err)
			//}
			go WaitAndReboot()
			record.UpgradeTime = time.Now().Unix()
			record.RebootPushed = true
		}
	}

	//if t.NeedReboot && record.RebootPushed && current_agent_version == record.AgentVersion {
	//
	//}

	err = fileutil.WriteToFileAsJson(f, record, "  ", true)
	if err != nil {
		logger.NotificationLogger().Warnf("failed WriteToFileAsJson after pushed, record:%+v, err:%v", record, err)
	} else {
		logger.NotificationLogger().Debugf("succ WriteToFileAsJson after pushed, record:%+v, f:%v", record, f)
	}

	RecheckUpgradeStatus(record.Rebooted)

	if hardware.RunningInDocker() {
		// 删除 aospace-upgrade 镜像
		dockerApi = dockerfacade.NewDockerFacade()
		containerInfos, err := dockerApi.ListContainers()
		if err != nil {
			return
		}
		for _, containerInfo := range containerInfos {
			if strings.Contains(containerInfo.Image, "space-upgrade") {
				err = dockerApi.RemoveContainer(containerInfo.ID)
				if err != nil {
					return
				}
				err = dockerApi.RemoveImage(containerInfo.ImageID)
				if err != nil {
					logger.UpgradeLogger().Errorf("remove upgrade image %s err:%v", containerInfo.ImageID, err)
					return
				}
			}
		}
	}

}
