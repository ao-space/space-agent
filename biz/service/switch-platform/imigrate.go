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
	"agent/config"
	"agent/utils"
	"fmt"
	"net/http"
	"strings"

	"agent/utils/logger"
	utilshttp "agent/utils/network/http"
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

func imigrate() (*ImigrateRsp, error) {

	var mi MigrateInfo
	mi.NetworkClinetId = device.GetDeviceInfo().NetworkClient.ClientID

	mi.UserInfos = si.OldAccount
	url, _ := utils.JoinUrl(si.NewApiBaseUrl, strings.ReplaceAll(config.Config.Platform.Migration.Path, "{box_uuid}", device.GetDeviceInfo().BoxUuid))

	destBRK, err := pair.GetDeviceRegKey(si.NewApiBaseUrl)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"Box-Reg-Key": destBRK.BoxRegKey, "Request-Id": si.TransId}
	rsp := ImigrateRsp{}

	httpReq, httpRsp, body, err1 := utilshttp.PostJsonWithHeaders(url, mi, headers, &rsp)
	if err1 != nil {
		logger.AppLogger().Warnf("Failed PostJson, transId:%v, err:%v, @@httpReq:%+v, @@httpRsp:%+v, @@body:%v", si.TransId, err1, httpReq, httpRsp, string(body))
		return nil, err1
	}
	logger.AppLogger().Infof("transId:%v, imigrate, parms:%+v", si.TransId, mi)
	logger.AppLogger().Infof("transId:%v, imigrate, rsp:%+v", si.TransId, rsp)
	logger.AppLogger().Infof("transId:%v, imigrate, httpReq:%+v", si.TransId, httpReq)
	logger.AppLogger().Infof("transId:%v, imigrate, httpRsp:%+v", si.TransId, httpRsp)
	logger.AppLogger().Infof("transId:%v, imigrate, body:%v", si.TransId, string(body))

	if httpRsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code:%v", httpRsp.StatusCode)
	}

	return &rsp, nil
}
