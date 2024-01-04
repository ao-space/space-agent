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
 * @Date: 2021-12-03 16:47:06
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-21 13:51:20
 * @Description:
 */

package pair

import (
	"fmt"
	"time"

	"agent/biz/docker"
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/dto"
	dtopair "agent/biz/model/dto/pair"
	"agent/biz/service/call"
	"agent/biz/service/encwrapper"
	"agent/config"
	"agent/utils"

	"agent/utils/logger"
)

// app与盒子的配对和初始化v1
func ServicePairing(req *dtopair.PairingReq) (dto.BaseRspStr, error) {

	err := encwrapper.Check()
	if err != nil {
		logger.AppLogger().Warnf("check, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}

	rt, err := encwrapper.Dec(req.ClientUuid, req.ClientPhoneModel)
	if err != nil {
		logger.AppLogger().Warnf("dec, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}
	clientUuid := rt[0]
	clientPhoneModel := rt[1]

	logger.AppLogger().Debugf("clientUuid:%+v, clientPhoneModel:%+v",
		clientUuid, clientPhoneModel)

	return doPairing(clientUuid, clientPhoneModel)
}

func doPairing(clientUuid, clientPhoneModel string) (dto.BaseRspStr, error) {

	err := ServiceRegisterBox()
	if err != nil {
		host, _, _ := utils.ParseUrl(device.GetApiBaseUrl())
		hostOfficial, _, _ := utils.ParseUrl(config.Config.Platform.APIBase.Url)
		if host != hostOfficial {
			err1 := fmt.Errorf("ServiceRegisterBox failed, current space platform:%v, %+v", host, err)
			logger.AppLogger().Warnf("%v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodePrivateSSPRegBoxErr, Message: err1.Error()}, err
		} else {
			err1 := fmt.Errorf("ServiceRegisterBox failed, %+v", err)
			logger.AppLogger().Warnf("%v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err1.Error()}, err
		}
	}
	device.SetBoxRegistered()

	clientinfo.SaveClientExchangeKey() // https://pm.eulix.xyz/bug-view-645.html

	// 盒子注册成功了. 开始启动容器. 后面将调用容器的注册管理员接口.
	go docker.PostEvent(docker.EventPairing)

	type CreateStruct struct {
		BoxUUID    string `json:"boxUUID,omitempty"`
		BoxRegKey  string `json:"boxRegKey,omitempty"`
		ClientUUID string `json:"clientUUID,omitempty"`
		PhoneModel string `json:"phoneModel,omitempty"`
		ApplyEmail string `json:"applyEmail,omitempty"`
	}
	req := &CreateStruct{BoxUUID: device.GetDeviceInfo().BoxUuid,
		BoxRegKey:  device.GetDeviceInfo().BoxRegKey,
		ClientUUID: clientUuid,
		PhoneModel: clientPhoneModel}

	applyEmail, _ := device.GetApplyEmail()
	if len(applyEmail) > 0 {
		req.ApplyEmail = applyEmail
	}

	var results call.MicroServerRsp
	err = call.CallServiceByPost(config.Config.Account.AdminCreate.Url, nil, req, &results)
	if err != nil {
		var err1 error
		for i := 0; i < 15; i++ {
			time.Sleep(time.Duration(4) * time.Second)

			err1 = call.CallServiceByPost(config.Config.Account.AdminCreate.Url, nil, req, &results)
			if err1 == nil {
				logger.AppLogger().Debugf("pair, CallServiceByPost return no err, accout return results:%+v",
					results)
				break
			}
			logger.AppLogger().Debugf("pair, waiting docker service started, err1:%v", err1)
		}
		logger.AppLogger().Debugf("pair, loop break, err1:%v", err1)

		if err1 != nil {
			return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err1.Error()},
				err
		}
	}
	return encwrapper.Enc(results)
}
