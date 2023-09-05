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

package notification

import (
	upModel "agent/biz/model/upgrade"
	"agent/utils/logger"
	"time"
)

const (
	DownloadedEvent = "upgrade_download_success"
	InstallingEvent = "upgrade_installing"
	SuccessEvent    = "upgrade_success"
	RestartEvent    = "upgrade_restart"
)

// OnUpgradeDownloadedSuccess push download successfully msg to redis by storeIntoRedis.
func OnUpgradeDownloadedSuccess(newVersion string) error {
	logger.NotificationLogger().Debugf("onUpgradeDownloadedSuccess")

	clientUUID, err := clientUuid()
	if err != nil {
		return err
	}

	// 开启自动安装的情况下, 下载成功不推送
	upgradeConf := upModel.GetUpgradeSettings()
	logger.NotificationLogger().Debugf("onUpgradeDownloadedSuccess, upgradeConf:%+v", upgradeConf)
	if upgradeConf.AutoInstall {
		return nil
	}
	type VersionInfo struct {
		Version string `json:"version"`
	}
	ver := &VersionInfo{Version: newVersion}
	var err1 error
	for i := 0; i < 3; i++ {
		_, err1 = storeIntoRedis(clientUUID, DownloadedEvent, ver)
		if err1 == nil {
			break
		}
		logger.NotificationLogger().Debugf("storeIntoRedis, waiting storeIntoRedis, err1:%v", err1)
		time.Sleep(time.Duration(1) * time.Second)
	}
	logger.NotificationLogger().Debugf("storeIntoRedis, loop break, err1:%v", err1)

	if err1 != nil {
		return err1
	}
	return nil
}

// OnUpgradeInstalling push installing msg to redis stream
func OnUpgradeInstalling() (string, error) {
	logger.NotificationLogger().Debugf("onUpgradeInstalling")

	clientUUID, err := clientUuid()
	if err != nil {
		return "", err
	}

	// 容器可能正在重启, 发失败了就不重试了. 要不然可能已经升级完成才发出去.
	id, err := storeIntoRedis(clientUUID, InstallingEvent, "")
	if err != nil {
		logger.NotificationLogger().Debugf("storeIntoRedis, waiting storeIntoRedis, err:%v", err)
		return "", err
	}
	return id, nil
}

// OnUpgradeSuccess push upgrade success msg to redis stream
func OnUpgradeSuccess() error {
	logger.NotificationLogger().Debugf("onUpgradeSuccess")

	clientUUID, err := clientUuid()
	if err != nil {
		return err
	}

	var err1 error
	for i := 0; i < 60; i++ {
		_, err1 := storeIntoRedis(clientUUID, SuccessEvent, "")
		if err1 == nil {
			break
		}
		logger.NotificationLogger().Debugf("storeIntoRedis, waiting storeIntoRedis, err1:%v", err1)
		time.Sleep(time.Duration(10) * time.Second)
	}
	logger.NotificationLogger().Debugf("storeIntoRedis, loop break, err1:%v", err1)

	if err1 != nil {
		return err1
	}

	return nil
}

// OnUpgradeRestart push upgrade restart to redis stream
func OnUpgradeRestart() error {
	logger.NotificationLogger().Debugf("onUpgradeRestart")
	clientUUID, err := clientUuid()
	if err != nil {
		return err
	}

	var err1 error
	for i := 0; i < 60; i++ {
		_, err1 := storeIntoRedis(clientUUID, RestartEvent, "")
		if err1 == nil {
			break
		}
		logger.NotificationLogger().Debugf("storeIntoRedis, waiting storeIntoRedis, err1:%v", err1)
		time.Sleep(time.Duration(10) * time.Second)
	}
	logger.NotificationLogger().Debugf("storeIntoRedis, loop break, err1:%v", err1)

	if err1 != nil {
		return err1
	}

	return nil
}
