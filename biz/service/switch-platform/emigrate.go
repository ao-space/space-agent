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
	"agent/biz/model/device"
	"agent/biz/service/pair"
	"agent/config"
	"agent/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"agent/utils/logger"
	utilshttp "agent/utils/network/http"
)

func emigrate() {
	type UserDomainRoute struct {
		UserId             string `json:"userId"`
		UserDomainRedirect string `json:"userDomainRedirect"`
	}
	type RouteReq struct {
		UserDomainRoutes []UserDomainRoute `json:"userDomainRouteInfos"`
	}

	type RouteRsp struct {
		BoxUUID          string            `json:"boxUUID"`
		UserDomainRoutes []UserDomainRoute `json:"userDomainRouteInfos"`
	}

	var params RouteReq
	for _, oldAC := range si.OldAccount {
		for _, newAC := range si.ImigrateResult.UserInfos {
			if oldAC.UserId == newAC.UserId {
				params.UserDomainRoutes = append(params.UserDomainRoutes, UserDomainRoute{UserId: oldAC.UserId, UserDomainRedirect: newAC.UserDomain})
			}
		}
	}

	//执行迁出
	for i := 0; i < 10; i++ {

		time.Sleep(time.Second * 5)

		url, _ := utils.JoinUrl(si.OldApiBaseUrl, strings.ReplaceAll(config.Config.Platform.Route.Path, "{box_uuid}", device.GetDeviceInfo().BoxUuid))

		destBRK, err := pair.GetDeviceRegKey(si.OldApiBaseUrl)
		if err != nil {
			logger.AppLogger().Infof("transId:%v,Failed to get box-reg-key", si.TransId)
			continue
		}

		headers := map[string]string{"Box-Reg-Key": destBRK.BoxRegKey, "Request-Id": si.TransId}
		rsp := RouteRsp{}

		httpReq, httpRsp, body, err := utilshttp.PostJsonWithHeaders(url, params, headers, &rsp)
		if err != nil {
			logger.AppLogger().Warnf("Failed PostJson, transId:%v,  err:%v, @@httpReq:%+v, @@httpRsp:%+v, @@body:%v", si.TransId, err, httpReq, httpRsp, string(body))
			continue
		}
		logger.AppLogger().Infof("transId:%v,emigrate, parms:%+v", si.TransId, params)
		logger.AppLogger().Infof("transId:%v,emigrate, rsp:%+v", si.TransId, rsp)
		logger.AppLogger().Infof("transId:%v,emigrate, httpReq:%+v", si.TransId, httpReq)
		logger.AppLogger().Infof("transId:%v,emigrate, httpRsp:%+v", si.TransId, httpRsp)
		logger.AppLogger().Infof("transId:%v,emigrate, body:%v", si.TransId, string(body))

		if httpRsp.StatusCode != http.StatusOK {
			logger.AppLogger().Infof("transId:%v, Failed to route. httpRsp.StatusCode:%v", si.TransId, httpRsp.StatusCode)
			UpdateStatus(StatusAbort, fmt.Sprintf("url:%v, StatusCode:%v", url, httpRsp.StatusCode))
		} else {
			UpdateStatus(StatusOK, "OK")
		}

		return
	}

	UpdateStatus(StatusAbort, "failed to emigrate, retry too much times.")
}
