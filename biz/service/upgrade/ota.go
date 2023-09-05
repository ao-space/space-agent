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
	"agent/biz/model/device_ability"
	"agent/biz/model/upgrade"
	"agent/utils/logger"
	"agent/utils/tools"
	"io"
	"net/http"
	"os"
	"time"
)

const OTAImagePathXZ = "/home/eulixspace_link/update.img.xz"
const OTAImagePath = "/home/eulixspace_link/update.img"

func DownloadOTAImgFile(url string) (upgrade.VersionDownInfo, error) {
	logger.UpgradeLogger().Infof("start to download ota kernel image from %s", url)
	file, err := os.OpenFile(OTAImagePathXZ, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return upgrade.VersionDownInfo{}, err
	}

	defer func() {
		_ = file.Close()
	}()

	rsp, err := http.Get(url)

	defer func() {
		rsp.Body.Close()
		file.Close()
	}()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return upgrade.VersionDownInfo{}, err
	}
	return upgrade.VersionDownInfo{
		Downloaded: true,
		UpdateTime: time.Now()}, nil
}

func OTAKernelUpgrade() {
	// 解压
	if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_SUPPORTED_FROM_MODEL_NUMBER {
		err := Unxz(OTAImagePathXZ)
		if err != nil {
			return
		}
		_, stdout, err := tools.RunCmd("updateEngine", []string{"--image_url=/home/eulixspace_link/update.img", "--misc=update", "--savepath=/home/eulixspace_link/update.img"})
		if err != nil {
			logger.UpgradeLogger().Errorf("OTA updrade error: %v", err)
			logger.UpgradeLogger().Debugf(stdout)
			return
		}

		RemoveOldImage()
	}
	return
}

func DnfKernelUpgrade(versionId string, kernelVersion string) error {
	if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_SUPPORTED_FROM_MODEL_NUMBER {
		// 4.19.90-2201.1.2.24
		_, stdout, err := tools.RunCmd("dnf", []string{"update", "-y", "rk3568-kernel-" + kernelVersion})
		if err != nil {
			logger.UpgradeLogger().Errorf("OTA updrade error: %v", err)
			db.MarkTaskInstallErr(versionId)
			return err
		}
		logger.UpgradeLogger().Debugf(stdout)
	}
	return nil
}

// RemoveOldImage 删除已安装的OTA镜像包
func RemoveOldImage() {
	if _, err := os.Stat(OTAImagePath); err == nil {
		os.Remove(OTAImagePath)
	}
	if _, err := os.Stat(OTAImagePathXZ); err == nil {
		os.Remove(OTAImagePathXZ)
	}
}

func GetCurrentKernelVersion() {

}

func Md5sum(path string) string {
	_, stdout, err := tools.RunCmd("md5sum", []string{path})
	if err != nil {
		logger.UpgradeLogger().Errorf("md5sum %s error: %v", path, err)
		logger.UpgradeLogger().Debugf(stdout)
		return ""
	}
	return stdout
}

func Unxz(path string) error {
	logger.UpgradeLogger().Infof("start to unxz %s", path)
	_, stdout, err := tools.RunCmd("unxz", []string{"-v", path})
	if err != nil {
		logger.UpgradeLogger().Errorf("unxz %s error: %v", path, err)
		logger.UpgradeLogger().Debugf(stdout)
		return err
	}
	return nil
}
