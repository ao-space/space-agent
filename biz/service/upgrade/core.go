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
	"agent/biz/db"
	"agent/biz/docker"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/upgrade"
	"agent/biz/notification"
	"agent/biz/service/base"
	"agent/biz/service/call"
	"agent/config"
	"agent/utils"
	"agent/utils/logger"
	"agent/utils/tools"
	"agent/utils/version"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/resty.v1"
)

type AgentUpgradeSvc struct {
	base.BaseService
}

func (svc *AgentUpgradeSvc) Process() {

}

type NativeAgent struct {
	VersionId string
}

type ContainerAgent struct {
	VersionId   string
	ComposeFile string
}

type Agent interface {
	Upgrade() error
	Download() error
	CurrentVersion() string
	AutoUpgrade()
}

func (na *NativeAgent) Upgrade() error {
	return nil
}

func (na *NativeAgent) Download() error {
	return nil
}

func (na *NativeAgent) CurrentVersion() string {
	return ""
}

func (ca *ContainerAgent) Upgrade() error {
	upgradeReq := upgrade.AllInOneUpgradeReq{
		VersionId: ca.VersionId,
		DataDir:   os.Getenv("AOSPACE_DATADIR"),
	}
	var microServerRsp call.MicroServerRsp
	err := call.CallServiceByPost(config.Config.Upgrade.Url, nil, &upgradeReq, &microServerRsp)
	if err != nil {
		db.MarkTaskInstallErr(ca.VersionId)
		logger.AppLogger().Errorf("upgrade CallServiceByPost:%v", err)
		return err
	}
	_, err = notification.OnUpgradeInstalling()
	if err != nil {
		logger.UpgradeLogger().Errorf("send installing notification failed,err:%v", err)
	}
	return nil
}

func (ca *ContainerAgent) Download() error {
	// start run upgrade container
	var agentRpmInfo = upgrade.VersionDownInfo{UpdateTime: time.Now()}
	var imageInfo = upgrade.VersionDownInfo{UpdateTime: time.Now()}
	var kernelInfo = upgrade.VersionDownInfo{UpdateTime: time.Now()}
	err := docker.ContainersUpAndPrune(config.Config.Docker.UpgradeComposeFile, nil)
	if err != nil {
		db.MarkTaskDownErr(ca.VersionId)
		return err
	}
	// start download other images
	imageInfo, err = PullImageFromCompose(ca.VersionId, ca.ComposeFile)
	if err != nil {
		db.MarkTaskDownErr(ca.VersionId)
		return fmt.Errorf("pull docker image %v, %v: %v", ca.VersionId, ca.ComposeFile, err)
	} else {
		db.MarkTaskDownloaded(ca.VersionId, agentRpmInfo, imageInfo, kernelInfo)
		err = notification.OnUpgradeDownloadedSuccess(ca.VersionId)
		if err != nil {
			logger.UpgradeLogger().Errorf("send download success notification failed,err:%v", err)
		}
		return nil
	}
}

func (ca *ContainerAgent) CurrentVersion() string {
	return config.VersionNumber
}

// GetLatestVersionMetadata is request platform for new version info and mini pkg
func GetLatestVersionMetadata() (upgrade.OverallInfo, error) {
	vI := upgrade.OverallInfo{}
	lastVersion, err := CheckLatestVersion()
	vI.VersionId = lastVersion.PkgVersion
	logger.UpgradeLogger().Debugf("lastVersion:%v", lastVersion)
	currentVer := config.VersionNumber
	if vI.VersionId == currentVer {
		return vI, dto.AlreadyLatestVersion
	}
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to check last compose: %s", err)
		return vI, err
	}
	downPath, err := DownloadComposeFile(lastVersion)
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to download compose file: %s", err)
		return vI, err
	}
	vI.Downloaded = true
	vI.PkgPath = downPath
	vI.Restart = lastVersion.Restart
	vI.KernelInfo.KernelVersion = lastVersion.KernelVersion
	vI.KernelInfo.KernelUrl = lastVersion.KernelUrl
	vI.KernelInfo.KernelMd5 = lastVersion.KernelMd5

	return vI, nil
}

func CheckLatestVersion() (upgrade.VersionFromPlatformV2, error) {
	apiBase := config.Config.Platform.APIBase.Url
	urlPath := config.Config.Platform.LatestVersionV2.Path
	versionDesc := upgrade.VersionFromPlatformV2{}

	checkUrl, err := utils.JoinUrl(apiBase, urlPath)
	if err != nil {
		return versionDesc, err
	}
	rId := uuid.New().String()

	_, err = resty.R().SetQueryParams(
		map[string]string{
			"box_pkg_name": "space.ao.server",
			"box_pkg_type": "box",
			"version_type": "open_source",
		}).SetHeader("Request-Id", rId).SetResult(&versionDesc).Get(checkUrl)

	if err != nil {
		return versionDesc, fmt.Errorf("%v: uuid: %v", err, rId)
	}

	return versionDesc, nil
}

// CheckVersionStatusRight is for preventing multiple execution on the same task
func CheckVersionStatusRight(versionId string) (int, error) {
	if versionId == "" {
		err := fmt.Errorf("param 'versionId' must be specified")
		return http.StatusBadRequest, err
	}
	task, err := db.ReadTask("")
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to get task record for %s: %s", versionId, err)
		return http.StatusInternalServerError, err
	}
	vRight := task.VersionId == versionId

	if vRight && task.Status == upgrade.Installing {
		err = fmt.Errorf("already exists a task is %s", upgrade.Installing)
		return http.StatusConflict, err
	} else {
		downloaded := vRight && task.DownStatus == upgrade.Done
		if !downloaded {
			err = fmt.Errorf("the version %s not been downloaded", versionId)
			return http.StatusBadRequest, err
		}
	}
	return 0, nil
}

func (na *NativeAgent) AutoUpgrade() {
	rand.Seed(time.Now().Unix())
	sec := rand.Intn(7200)
	sTime := time.Duration(sec) * time.Second
	future := time.Now().Add(sTime)
	logger.UpgradeLogger().Infof("[auto-upgrade] will be starting at %s", future.Format(time.RFC3339))
	time.After(sTime)
	logger.UpgradeLogger().Infof("[auto-upgrade] Start to check new version regularly")

	// upConf := config.Config.Box.UpgradeConfig
	upConf := upgrade.GetUpgradeSettings()

	// 根据产品要求不开启自动下载，无法开启自动安装
	if upConf.AutoDownload {
		versionInfo, err := GetLatestVersionMetadata()
		if err != nil {
			logger.UpgradeLogger().Errorf("[auto-upgrade] GetLatestVersionMetadata error : %v", err)
			return
		}
		//var kernelInfo *upgrade.VersionDownInfo

		if versionInfo.Restart {
			logger.UpgradeLogger().Infof("[auto-upgrade] this upgrade version need to restart device")
			db.MarkTaskRebootFlag(versionInfo.VersionId)
		}
		var kernelInfo = upgrade.VersionDownInfo{UpdateTime: time.Now()}
		if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_SUPPORTED_FROM_MODEL_NUMBER {
			// 内核OTA升级
			if versionInfo.KernelUrl != "" {
				kernelInfo, err = DownloadOTAImgFile(versionInfo.KernelUrl)
				if err != nil {
					logger.UpgradeLogger().Errorf("[auto-upgrade] download ota kernel image error:%v", err)
					return
				}
				localMd5 := Md5sum(OTAImagePathXZ)
				logger.UpgradeLogger().Infof("[auto-upgrade] start to verify md5")
				logger.UpgradeLogger().Debugf("[auto-upgrade] local md5:%s", localMd5)
				logger.UpgradeLogger().Debugf("[auto-upgrade] versionInfo.KernelMd5 :%s", versionInfo.KernelMd5)
				//校验MD5
				//校验MD5
				localMd5Output := Md5sum(OTAImagePathXZ)
				localMd5Split := strings.Split(localMd5Output, " ")
				if len(localMd5Split) > 1 {
					if versionInfo.KernelMd5 == localMd5Split[0] && versionInfo.KernelMd5 != "" {
						logger.UpgradeLogger().Infof("[auto-upgrade] md5sum verify passed")
						logger.UpgradeLogger().Infof("[auto-upgrade] start to upgrade kernel %s", versionInfo.KernelVersion)
						err = DnfKernelUpgrade(versionInfo.VersionId, versionInfo.KernelVersion)
						if err != nil {
							logger.UpgradeLogger().Errorf("dnf update kernel error:%v", err)
							return
						}
						OTAKernelUpgrade()
					} else {
						logger.UpgradeLogger().Infof("[auto-upgrade] md5sum verify failed")
						db.MarkTaskDownErr(versionInfo.VersionId)
						return
					}
				}
			}
		}

		// 下载system-agent RPM
		rpmInfo, err := DownloadRpm(versionInfo.VersionId, AgentName)
		if err != nil {
			db.MarkTaskDownErr(versionInfo.VersionId)
			logger.UpgradeLogger().Errorf("[auto-upgrade] download system-agent rpm error,%v", err)
			return
		}
		// 下载docker镜像

		imageInfo, err := PullImageFromCompose(versionInfo.VersionId, versionInfo.PkgPath)
		if err != nil {
			db.MarkTaskDownErr(versionInfo.VersionId)
			logger.UpgradeLogger().Errorf("[auto-upgrade] pull docker image error,%v", err)
			return
		} else {
			db.MarkTaskDownloaded(versionInfo.VersionId, rpmInfo, imageInfo, kernelInfo)
		}

		if upConf.AutoInstall && VersionDownloaded(versionInfo.VersionId) {
			task, err := db.MarkTaskInstalling(versionInfo.VersionId)
			if err != nil {
				logger.UpgradeLogger().Errorf("[auto-upgrade] Failed to mark task installing: %s", err)
			}
			// 升级并重启system-agent
			InstallAgentV2(task.VersionId)
		}
	}
}

func (ca *ContainerAgent) AutoUpgrade() {
	rand.Seed(time.Now().Unix())
	sec := rand.Intn(7200)
	sTime := time.Duration(sec) * time.Second
	future := time.Now().Add(sTime)
	logger.UpgradeLogger().Infof("[auto-upgrade] will be starting at %s", future.Format(time.RFC3339))
	time.After(sTime)
	logger.UpgradeLogger().Infof("[auto-upgrade] Start to check new version regularly")

	// upConf := config.Config.Box.UpgradeConfig
	upConf := upgrade.GetUpgradeSettings()

	// 根据产品要求不开启自动下载，无法开启自动安装
	if upConf.AutoDownload {
		versionInfo, err := GetLatestVersionMetadata()
		if err != nil {
			logger.UpgradeLogger().Errorf("[auto-upgrade] GetLatestVersionMetadata error : %v", err)
			return
		}
		ca.VersionId = versionInfo.VersionId
		ca.ComposeFile = versionInfo.PkgPath
		err = ca.Download()
		if err != nil {
			logger.UpgradeLogger().Errorf("[auto upgrade] download err:%v", err)
			return
		}
		if upConf.AutoInstall && VersionDownloaded(versionInfo.VersionId) {
			_, err := db.MarkTaskInstalling(versionInfo.VersionId)
			if err != nil {
				logger.UpgradeLogger().Errorf("[auto-upgrade] Failed to mark task installing: %s", err)
			}
			// 升级并重启aospace-all-in-one容器
			err = ca.Upgrade()
			if err != nil {
				logger.UpgradeLogger().Errorf("[auto upgrade] download err:%v", err)
				return
			}
		}
	}

}

// 主要是防止升级过程中被外部原因意外终止，一直停留在 installing 状态
func RecheckUpgradeStatus(rebooted bool) {
	// check the status
	logger.UpgradeLogger().Debugf("Rechecking the status of upgrade after started")
	task, err := db.ReadTask("")
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to recheck task status: %v", err)
		return
	}
	// change the status
	versionId := version.GetInstalledAgentVersionRemovedNewLine()
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to recheck upgrade Status %v", err)
		return
	}
	if len(versionId) < 3 {
		versionId = config.VersionNumber
	}
	needAmend := false
	logger.UpgradeLogger().Debugf("print task: %v", task)
	if task.Status == upgrade.Installing {
		logger.UpgradeLogger().Warnf("Find have installed version is '%s', the upgrade status is '%s': '%s',that may be not right,"+
			" will be to amend it", versionId, task.VersionId, task.Status)
		needAmend = true
		logger.UpgradeLogger().Debugf("current version:%s", versionId)
		logger.UpgradeLogger().Debugf("task version:%s", task.VersionId)
		logger.UpgradeLogger().Debugf("version compare:%v", task.VersionId == versionId)

		if versionId != task.VersionId {
			task.Status = upgrade.InstallErr
			task.InstallStatus = upgrade.Err
			task.DoneInstallTime = time.Now().Format(time.RFC3339)
		} else if versionId == task.VersionId && !task.NeedReboot {
			task.Status = upgrade.Installed
			task.InstallStatus = upgrade.Done
			task.DoneInstallTime = time.Now().Format(time.RFC3339)
		} else if versionId == task.VersionId && rebooted && task.NeedReboot {
			task.Status = upgrade.Installed
			task.InstallStatus = upgrade.Done
			task.DoneInstallTime = time.Now().Format(time.RFC3339)
		}
	} else if task.Status == upgrade.Downloading {
		logger.UpgradeLogger().Warnf("Find task status is %s that may be not right, "+
			"will be change to %s", upgrade.Downloading, upgrade.DownloadErr)
		task.Status = upgrade.DownloadErr
		task.DownStatus = upgrade.Err
		task.DoneDownTime = time.Now().Format(time.RFC3339)
		needAmend = true
	} else if task.Status == upgrade.Downloaded && task.VersionId != versionId {
		task.VersionId = versionId
		needAmend = true
		return
	}
	if needAmend {
		logger.UpgradeLogger().Debugf("print task :%v", task)
		doc, err := db.UpdateOrCreateTask(task)
		if err != nil {
			logger.UpgradeLogger().Errorf("Failed to chang upgrade status %s", err)
			return
		}
		logger.UpgradeLogger().Infof("upgrade status had changed: %v", doc)
	}
	if task.Status == upgrade.Installed {
		RemoveOldRpmPkg(task.VersionId)
	}
	logger.UpgradeLogger().Infof("Recheck status done")
}

func VersionDownloaded(versionId string) bool {
	task, err := db.ReadTask(versionId)
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to check version %s if Downloaded: %s", versionId, err)
		return false
	}
	return task.DownStatus == upgrade.Done
}

func RemoveOldRpmPkg(versionId string) {
	task, err := db.ReadTask(versionId)
	if err != nil {
		logger.UpgradeLogger().Errorf("read task error:%v", err)
	}
	// 删除旧的system-agent rpm包
	if _, err = os.Stat(task.RpmPkg.PkgPath); err == nil {
		logger.UpgradeLogger().Infof("remove old rpm pkg :%s", task.RpmPkg.PkgPath)
		err = os.Remove(task.RpmPkg.PkgPath)
		if err != nil {
			logger.UpgradeLogger().Errorf("read task error:%v", err)
			return
		}
	}
}

func WaitAndReboot() {
	logger.UpgradeLogger().Debugf("device will reboot after 90 seconds")
	duration := 90 * time.Second
	ticker := time.Tick(time.Second)

	for remaining := duration; remaining > 0; remaining -= time.Second {
		logger.UpgradeLogger().Debugf("device will reboot after %d seconds ...\n", remaining/time.Second)
		<-ticker
	}
	tools.ExeCmd("reboot")
}
