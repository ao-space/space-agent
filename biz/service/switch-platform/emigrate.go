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
 * @Date: 2022-11-15 14:53:56
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-11-24 17:30:04
 */

package switchplatform

import (
	"agent/biz/service/pair"
	"agent/utils/logger"
	"github.com/big-dust/platform-sdk-go/v2"
	"time"
)

func emigrate() {

	var params platform.SpacePlatformMigrationOutRequest
	for _, oldAC := range si.OldAccount {
		for _, newAC := range si.ImigrateResult.UserInfos {
			if oldAC.UserId == newAC.UserId {
				params.UserDomainRouteInfos = append(params.UserDomainRouteInfos, platform.UserDomainRouteInfo{UserId: oldAC.UserId, UserDomainRedirect: newAC.UserDomain})
			}
		}
	}

	//执行迁出
	for i := 0; i < 10; i++ {

		time.Sleep(time.Second * 5)

		client, err := pair.GetSdkClientWithDeviceRegKey(si.OldApiBaseUrl)
		if err != nil {
			logger.AppLogger().Warnf("Failed new SDK Client: err: %v", err.Error())
			return
		}
		resp, err := client.SpacePlatformMigrationOut(&params)

		if err != nil {
			logger.AppLogger().Warnf("Space Platform Migrate Out Failed by SDK: transId:%v,emigrate, err:%+v ", si.TransId, err)
			continue
		}
		logger.AppLogger().Infof("transId:%v,emigrate, parms:%+v", si.TransId, params)
		logger.AppLogger().Infof("transId:%v,emigrate, rsp:%+v", si.TransId, resp)

		return
	}

	UpdateStatus(StatusAbort, "failed to emigrate, retry too much times.")
}
