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
 * @Author: zhongguang
 * @Date: 2022-12-12 11:46:05
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-12-12 13:35:35
 */

package notification

import (
	"agent/biz/model/device"
	"agent/utils"
	"sync"
	"time"

	"agent/utils/logger"
)

var lastVersion string
var mtx sync.Mutex

func DealVerisonChange(version string) {
	host, _, _ := utils.ParseUrl(device.GetApiBaseUrl())
	verInfo := host + ":" + version

	mtx.Lock()
	defer mtx.Unlock()

	if verInfo != lastVersion && onAbilityChange() == nil {
		lastVersion = verInfo
	}
}

func onAbilityChange() error {
	logger.NotificationLogger().Debugf("onAbilityChange")

	optType := "ability_change"
	clientUUID, err := clientUuid()
	if err != nil {
		return err
	}

	var err1 error
	for i := 0; i < 60; i++ {
		_, err1 := storeIntoRedis(clientUUID, optType, "")
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
