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
	"agent/biz/model/upgrade"
	"agent/config"
	"agent/utils/docker/dockerfacade"
	"agent/utils/hardware"
	"agent/utils/logger"
	"agent/utils/tools"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

const (
	OsType      = "aarch64"
	AgentName   = "eulixspace-agent"
	UpgradeName = "eulixspace-upgrade"
)

var dockerApi *dockerfacade.DockerFacade

//func DownloadPkgQuiet(versionId string, cFile string, kernel upgrade.KernelInfo) {
//	// 获取最新的固件RPM包名和版本号
//	firmwares := upgrade.LatestUpgradeFirmware
//	logger.UpgradeLogger().Debugf("firmwares:%v", firmwares)
//	if firmwares != nil {
//		err := DownloadPkgs(versionId, cFile, kernel)
//		if err != nil {
//			logger.UpgradeLogger().Errorf("DownloadPkgQuiet: %s", err)
//		}
//	}
//
//}

func DownloadPkgs(versionId string, cFile string, kernel upgrade.KernelInfo) error {
	logger.UpgradeLogger().Debugf("Start to download pkg for version " + versionId)

	var agentRpmInfo = upgrade.VersionDownInfo{UpdateTime: time.Now()}
	var imageInfo = upgrade.VersionDownInfo{UpdateTime: time.Now()}
	var kernelInfo = upgrade.VersionDownInfo{UpdateTime: time.Now()}
	if hardware.RunningInDocker() {
		// start upgrade container
		err := docker.ContainersUpAndPrune(config.Config.Docker.UpgradeComposeFile, nil)
		if err != nil {
			db.MarkTaskDownErr(versionId)
			return err
		}
	} else {
		err := CleanAllAndMakeCache()
		if err != nil {
			logger.UpgradeLogger().Errorf("Failed to clean and make cache to dnf %v", err)
		}
		// 下载system-agent RPM 包
		agentRpmInfo, err = DownloadRpm(versionId, AgentName)
		if err != nil {
			db.MarkTaskDownErr(versionId)
			logger.UpgradeLogger().Errorf("download %s-%s error:%v", AgentName, versionId, err)
			return err
		} else {
			logger.UpgradeLogger().Debugf("download system-agent rpm successfully")
		}
	}
	// 下载OTA内核升级包，只有210 odm 盒子才能ota升级内核
	//var kernelDownInfo *upgrade.VersionDownInfo

	if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_SUPPORTED_FROM_MODEL_NUMBER {
		if _, err := os.Stat(OTAImagePathXZ); err != nil {
			if kernel.KernelUrl != "" {
				kernelInfo, err = DownloadOTAImgFile(kernel.KernelUrl)
				if err != nil {
					logger.UpgradeLogger().Errorf("download ota kernel image error:%v", err)
					return err
				}
				//校验MD5
				localMd5Output := Md5sum(OTAImagePathXZ)
				localMd5Split := strings.Split(localMd5Output, " ")
				if len(localMd5Split) > 1 {
					if kernel.KernelMd5 != localMd5Split[0] {
						logger.UpgradeLogger().Errorf("verify md5 failed,kernal md5:%s, %s md5 : %s", kernel.KernelMd5, OTAImagePathXZ, localMd5Split[0])
						db.MarkTaskDownErr(versionId)
						return err
					}
				}
			}
		} else {
			//校验MD5
			localMd5Output := Md5sum(OTAImagePathXZ)
			localMd5Split := strings.Split(localMd5Output, " ")
			if len(localMd5Split) > 1 {
				if kernel.KernelMd5 != localMd5Split[0] {
					logger.UpgradeLogger().Errorf("verify md5 failed,kernal md5:%s, %s md5 : %s", kernel.KernelMd5, OTAImagePathXZ, localMd5Split[0])
					db.MarkTaskDownErr(versionId)
					return err
				}
			}
		}
	}
	imageInfo, err := PullImageFromCompose(versionId, cFile)
	if err != nil {
		db.MarkTaskDownErr(versionId)
		return fmt.Errorf("pull docker image %v, %v: %v", versionId, cFile, err)
	} else {
		db.MarkTaskDownloaded(versionId, agentRpmInfo, imageInfo, kernelInfo)
		return nil
	}
}

func DownloadRpm(versionId string, rpmName string) (upgrade.VersionDownInfo, error) {
	logger.UpgradeLogger().Debugf("start to download rpm: %s", rpmName)
	pkgName := BuildPkgName(rpmName, versionId)
	downInfo := upgrade.VersionDownInfo{VersionId: versionId, Downloaded: false}
	saveDir := filepath.Join(config.Config.RunTime.BasePath, config.Config.RunTime.PkgDir)
	if strings.Contains(versionId, " ") {
		return downInfo, fmt.Errorf("version id cannot be contains blank")
	}
	_, stdout, err := tools.RunCmd("dnf", []string{"download", "--downloaddir=" + saveDir, pkgName, "-y"})
	if err != nil {
		return downInfo, fmt.Errorf("error downloading %s rpm %s :%s", pkgName, err, stdout)
	}
	downInfo.Downloaded = true
	downInfo.PkgPath = path.Join(saveDir, pkgName+"."+OsType+".rpm")
	downInfo.UpdateTime = time.Now()
	return downInfo, nil
}

func BuildPkgName(rpmName string, versionId string) string {
	return rpmName + "-" + versionId
}

func PullImageFromCompose(versionId string, cFile string) (upgrade.VersionDownInfo, error) {
	downInfo := upgrade.VersionDownInfo{VersionId: versionId}
	//docker.LoginRegistry()
	err := docker.ProcessEnv(cFile, nil)
	if err != nil {
		return downInfo, fmt.Errorf("PullImageFromCompose: %s", err)
	}
	logger.UpgradeLogger().Debugf("start to pull docker images")
	err = dockerApi.Pull(cFile)
	if err != nil {
		return downInfo, fmt.Errorf("PullImageFromCompose: %s", err)
	}

	downInfo.Downloaded = true
	downInfo.UpdateTime = time.Now()

	return downInfo, nil
}

func DownloadComposeFile(version upgrade.VersionFromPlatformV2) (string, error) {
	versionUrl := version.DownloadUrl
	logger.UpgradeLogger().Infof("Start to download last compose file...")
	savePath := filepath.Join(config.Config.RunTime.BasePath, config.Config.RunTime.PkgDir, "docker-compose.yml")
	logger.UpgradeLogger().Debugf("From url %s downloaded compose file to %s", versionUrl, savePath)
	err := DownFile(versionUrl, savePath)

	if err != nil {
		return savePath, err
	}

	return savePath, nil
}

func DownFile(url string, path string) error {
	// https://artifactory.eulix.xyz/artifactory/cicada-public/eulixspace-box/docker-compose-0.4.0-alpha.91033.yml
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = fileutil.WriteToFile(path, f, true)
	return nil
}
