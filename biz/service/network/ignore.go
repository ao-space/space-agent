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

package network

import (
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/dto/network"
	"agent/biz/service/base"
	"fmt"
	"time"

	utilsnetwork "agent/utils/rpi/network"

	"agent/utils/logger"
)

type NetworkIgnoreService struct {
	base.BaseService
}

func (svc *NetworkIgnoreService) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("NetworkIgnoreService")
	abilityModel := device_ability.GetAbilityModel()
	if !abilityModel.InnerDiskSupport {
		err := fmt.Errorf("unsupported function")
		return dto.BaseRspStr{Code: dto.AgentCodeUnsupportedFunction,
			Message: err.Error()}
	}

	req := svc.Req.(*network.NetworkIgnoreReq)
	// logger.AppLogger().Debugf("PostNetworkConfigService, req:%+v", req)

	succ, err := utilsnetwork.ForgetWifi(req.WIFIName)
	if err != nil {
		logger.AppLogger().Warnf("ForgetWifi err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err.Error()}
	}
	if !succ {
		err1 := fmt.Errorf("ForgetWifi succ==false")
		logger.AppLogger().Warnf("%v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err1.Error()}
	}
	time.Sleep(3 * time.Second) // 命令返回了，但是ip还在。时间关系，暂时这么处理。
	return svc.BaseService.Process()
}
