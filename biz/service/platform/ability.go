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

package platform

import (
	"agent/biz/model/device"
	"agent/biz/model/platform"
	"agent/biz/service/call"
	"agent/config"
	"agent/utils"
	"agent/utils/logger"
	"github.com/dungeonsnd/gocom/encrypt/random"
)

var platformApis *platform.PlatformAPIs

func InitPlatformAbility() *platform.PlatformAPIs {
	var headers = map[string]string{
		"Request-Id": random.GenUUID(),
	}
	apiBaseUrl := device.GetApiBaseUrl()
	logger.AppLogger().Debugf("apiBaseUrl:%s", apiBaseUrl)
	url, _ := utils.JoinUrl(apiBaseUrl, config.Config.Platform.Ability.Path)
	err := call.CallServiceByGet(url, headers, nil, &platformApis)
	if err != nil {
		logger.AppLogger().Errorf("Get Platform Ability Request error:%v", err)
		return nil
	}
	logger.AppLogger().Infof("Get Platform Ability Request Successfully")
	logger.AppLogger().Debugf("platformApis: %v", platformApis)
	return platformApis
}

func CheckPlatformAbility(uri string) bool {
	if platformApis == nil {
		platformApis = InitPlatformAbility()
	}

	for _, apis := range platformApis.PlatformAPIs {
		//logger.AppLogger().Debugf(apis.URI)
		if uri == apis.URI {
			return true
		}
	}
	return false

}
