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

package config

import (
	"agent/biz/model/device"
	"agent/biz/model/dto"
	"agent/biz/model/dto/bind/internet/service/config"
	"agent/biz/service/base"
	"agent/biz/service/call"
	"agent/biz/service/pair"
	cfg "agent/config"
	"agent/utils/logger"
	"fmt"
)

type InternetServiceGetConfig struct {
	base.BaseService
}

type UserInfo struct {
	ClientUuid string `json:"clientUUID" form:"clientUUID"`
	UserDomain string `json:"userDomain" form:"userDomain"`
	Aoid       string `gorm:"column:aoid" json:"aoId" form:"aoId"`
}

// GetUserInfo 从网关获取userinfo并序列化返回
func GetUserInfo() ([]*UserInfo, error) {
	type Rsp struct {
		Code      string      `json:"code"`
		RequestId string      `json:"requestId, omitempty"`
		Message   string      `json:"message, omitempty"`
		Results   []*UserInfo `json:"results, omitempty"`
	}

	var results Rsp
	url := cfg.Config.Account.Member.Url + "?userId=1"
	err := call.CallServiceByGet(url, nil, nil, &results)
	// logger.AppLogger().Debugf("InternetServiceGetConfig Process, err:%v, url:%v, results:%+v", err, url, results)
	if err != nil {
		return nil, err
	}

	return results.Results, nil
}

func (svc *InternetServiceGetConfig) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("InternetServiceGetConfig Process, svc.RequestId:%v", svc.RequestId)

	if svc.Req == nil {
		err := fmt.Errorf("req is nil")
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, RequestId: svc.RequestId, Message: err.Error()}
	}
	req := svc.Req.(*config.GetConfigReq)
	// logger.AppLogger().Debugf("InternetServiceGetConfig Process, req:%v", req)

	userInfos, err := GetUserInfo()
	if err != nil {
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, RequestId: svc.RequestId, Message: err.Error()}
	}
	userDomain := ""
	for _, userInfo := range userInfos {
		logger.AppLogger().Debugf("InternetServiceGetConfig Process, userInfo:%+v", userInfo)
		if len(req.Aoid) > 0 {
			if userInfo.Aoid == req.Aoid {
				userDomain = userInfo.UserDomain
			}
		} else if len(req.ClientUUID) > 0 {
			if userInfo.ClientUuid == req.ClientUUID {
				userDomain = userInfo.UserDomain
			}
		}
	}

	rsp := &config.GetConfigRsp{EnableInternetAccess: device.GetConfig().EnableInternetAccess,
		EnableP2P:       device.GetConfig().EnableInternetAccess,
		EnableLAN:       true,
		UserDomain:      userDomain,
		PlatformApiBase: device.GetApiBaseUrl(),
	}
	rsp.ConnectedNetwork = pair.GetConnectedNetwork()
	svc.Rsp = rsp
	return svc.BaseService.Process()
}
