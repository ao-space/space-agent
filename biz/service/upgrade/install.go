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
	"agent/biz/model/upgrade"
	"agent/biz/service/call"
	"agent/config"
	"agent/utils/hardware"
	"agent/utils/logger"
	"agent/utils/tools"
	"agent/utils/unixsock/http"
	"fmt"
	"os"
)

func dnfInstall(rpmPath string) error {
	logger.UpgradeLogger().Debugf("Start to install %s with dnf", rpmPath)
	_, stdErr, err := tools.ExeCmd("dnf", "update", rpmPath, "-y")
	if err != nil || stdErr != "" {
		return fmt.Errorf("dnf install: %s: %s", err, stdErr)
	}
	logger.UpgradeLogger().Infof("Success to install %s with dnf", rpmPath)
	return nil
}

func InstallAgent(versionId string) {
	logger.UpgradeLogger().Infof("install and restart system-agent,version:%s", versionId)
	//_, outMsg, _ := tools.RunCmd("nohup", []string{"eulixspace-upgrade", "install", "-v", versionId, ">/dev/null 2>&1 &"})
	logger.UpgradeLogger().Debugf(fmt.Sprintf("nohup eulixspace-upgrade install -v %s > /tmp/upgrade.log  2>&1 &", versionId))
	_, outMsg, err := tools.RunCmd("sh", []string{"/tmp/upgrade.sh", versionId})
	logger.UpgradeLogger().Debugf("out msg:%v", outMsg)
	if err != nil {
		logger.UpgradeLogger().Debugf("InstallAgent error: %v", err)
	}
	return
}

func InstallAgentV2(versionId string) {
	type Version struct {
		VersionId string
	}

	if hardware.RunningInDocker() {
		// 如果运行在容器中，调用upgrade的异步接口来执行all-in-one的升级和重启
		upgradeReq := upgrade.AllInOneUpgradeReq{
			VersionId: versionId,
			DataDir:   os.Getenv("AOSPACE_DATADIR"),
		}
		var microServerRsp call.MicroServerRsp
		err := call.CallServiceByPost(config.Config.Upgrade.Url, nil, &upgradeReq, &microServerRsp)
		if err != nil {
			logger.AppLogger().Errorf("upgrade CallServiceByPost:%v", err)
			return
		}
	} else {
		// 发送消息到aospace-upgrade ,由upgrade 升级并重启 system-agent
		err := http.SendMessageToSocket(config.Config.RunTime.BasePath+config.Config.RunTime.SocketFile, versionId)
		if err != nil {
			logger.UpgradeLogger().Errorf("send message to socket error:%v", err)
			db.MarkTaskInstallErr(versionId)
		}
	}

	return
}

//// InstallFirmwares 安装固件
//func InstallFirmwares(versionId string) error {
//	for _, firmware := range upgrade.LatestUpgradeFirmware.Firmwares {
//		pkgName := filepath.Join(config.Config.RunTime.BasePath+config.Config.RunTime.PkgDir, firmware.Package+"-"+firmware.Version+"."+OsType+".rpm")
//		if _, err := os.Stat(pkgName); err == nil {
//
//			err := dnfInstall(pkgName)
//			if err != nil {
//				logger.UpgradeLogger().Errorf("dnf install package %s error : %v", pkgName, err)
//				db.MarkTaskInstallErr(versionId)
//				return err
//			}
//			// 删除已安装的固件包
//			os.Remove(pkgName)
//		}
//	}
//	logger.UpgradeLogger().Debugf("no firmware need to update")
//	return nil
//}
