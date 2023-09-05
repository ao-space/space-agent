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
 * @Date: 2021-12-02 23:00:43
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-10 14:32:32
 * @Description:
 */

package upgrade

import (
	"agent/config"
	"agent/utils/file/storage"
	"agent/utils/logger"
)

var upSettingsStorage *storage.Storage

type UpgradeSettings struct {
	AutoDownload bool `json:"autoDownload"`
	AutoInstall  bool `json:"autoInstall"`
}

func init() {
	upSettingsStorage = storage.NewStorage(config.Config.Box.UpgradeConfig.SettingsFile)
}

func GetUpgradeSettings() *UpgradeSettings {
	settings := &UpgradeSettings{AutoDownload: true, AutoInstall: true}
	err := upSettingsStorage.LoadJson(settings)
	if err != nil {
		logger.AppLogger().Warnf("failed GetUpgradeSettings, file:%v, err:%v",
			config.Config.Box.UpgradeConfig.SettingsFile, err)
	} else {
		logger.AppLogger().Debugf("succ GetUpgradeSettings, file:%v, settings:%+v",
			config.Config.Box.UpgradeConfig.SettingsFile, settings)
	}
	return settings
}

func SetUpgradeSettings(settings *UpgradeSettings) error {
	err := upSettingsStorage.SaveJson(settings)
	if err != nil {
		logger.AppLogger().Warnf("failed SetUpgradeSettings, file:%v, err:%v",
			config.Config.Box.UpgradeConfig.SettingsFile, err)
		return err
	}
	logger.AppLogger().Debugf("succ SetUpgradeSettings, file:%v, settings:%+v",
		config.Config.Box.UpgradeConfig.SettingsFile, settings)
	return nil
}
