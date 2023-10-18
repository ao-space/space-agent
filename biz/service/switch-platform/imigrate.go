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
 * @Date: 2022-11-15 14:55:20
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-11-16 10:39:30
 */
package switchplatform

import (
	"agent/biz/model/device"
	"agent/biz/service/pair"
	"agent/utils/logger"
	"fmt"
	"github.com/big-dust/platform-sdk-go/v2"
	"github.com/jinzhu/copier"
)

type ImigrateRsp struct {
	NetworkClient device.NetworkClientInfo `json:"networkClient"`
	UserInfos     []AccountInfo            `json:"userInfos"`
}

func (r *ImigrateRsp) GetAdminDomain() (string, error) {
	for _, account := range r.UserInfos {
		if account.UserType == "user_admin" {
			return account.UserDomain, nil
		}
	}
	return "", fmt.Errorf("not found admin")
}

func imigrate() (*platform.SpacePlatformMigrationResponse, error) {

	var mi MigrateInfo
	mi.NetworkClinetId = device.GetDeviceInfo().NetworkClient.ClientID

	mi.UserInfos = si.OldAccount

	client, err := pair.GetSdkClientWithDeviceRegKey(si.NewApiBaseUrl)
	if err != nil {
		logger.AppLogger().Warnf("Failed new SDK Client: err: %v", err.Error())
		return nil, err
	}

	input := &platform.SpacePlatformMigrationRequest{
		NetworkClientId: mi.NetworkClinetId,
		UserInfos:       []platform.UserMigrationInfo{},
	}

	err = copier.Copy(&input.UserInfos, &mi.UserInfos)
	if err != nil {
		logger.AppLogger().Warnf("Copy UserInfos: err: %v", err.Error())
		return nil, err
	}
	resp, err := client.SpacePlatformMigration(input)

	if err != nil {
		logger.AppLogger().Warnf("Imgrate By Sdk Failed: %v", err.Error())
		return nil, err
	}

	logger.AppLogger().Infof("transId:%v, imigrate, rsp:%+v", si.TransId, resp)

	return resp, nil
}
