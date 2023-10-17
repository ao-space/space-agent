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
	"agent/biz/service/pair"
	"agent/utils/logger"
	"github.com/big-dust/platform-sdk-go/v2"
)

var platformApis *platform.GetAbilityResponse

func InitPlatformAbility() *platform.GetAbilityResponse {

	apiBaseUrl := device.GetApiBaseUrl()
	logger.AppLogger().Debugf("apiBaseUrl:%s", apiBaseUrl)

	client, err := pair.GetSdkClientWithDeviceRegKey(apiBaseUrl)
	if err != nil {
		logger.AppLogger().Errorf("Get SDK Client Request error:%v", err)
		return nil
	}
	platformApis, err = client.GetAbility()
	logger.AppLogger().Infof("Get Platform Ability Request Successfully")
	logger.AppLogger().Debugf("platformApis: %v", platformApis)
	return platformApis
}

func CheckPlatformAbility(uri string) bool {
	if platformApis == nil {
		platformApis = InitPlatformAbility()
	}

	for _, api := range platformApis.PlatformApis {
		//logger.AppLogger().Debugf(apis.URI)
		if uri == api.Uri {
			return true
		}
	}
	return false

}
