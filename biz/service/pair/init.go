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
 * @Date: 2021-12-13 14:55:34
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-28 10:58:10
 * @Description:
 */

package pair

import (
	"agent/biz/alivechecker/model"
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	dtodevice "agent/biz/model/dto/device"
	dtopair "agent/biz/model/dto/pair"
	"agent/biz/service/encwrapper"
	"agent/config"
	"fmt"

	"agent/utils/logger"
)

func ServiceInit() (dto.BaseRspStr, error) {
	logger.AppLogger().Debugf("ServiceInit")
	logger.AccessLogger().Debugf("[ServiceInit]")
	err := encwrapper.Check()
	if err != nil {
		err1 := fmt.Errorf("check failed, err:%v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err1.Error()}, err1
	}

	// https://code.eulix.xyz/bp/cicada/-/issues/465#note_82345
	boxUuid := device.GetDeviceInfo().BoxUuid

	paired := clientinfo.GetAdminPairedStatus()
	pairedBool := clientinfo.GetAdminPairedStatus() == clientinfo.DeviceAlreadyBound
	logger.AppLogger().Debugf("ServiceInit, paired:%v", paired)

	connected := 1
	if model.Get().PingCloudHost {
		logger.AppLogger().Debugf("ServiceInit, connected")
		connected = 0
	} else {
		logger.AppLogger().Debugf("ServiceInit, not connected")
	}

	results := &dtopair.InitResult{
		BoxName:                "傲空间",
		ProductId:              boxUuid,
		BoxUuid:                boxUuid,
		Paired:                 paired,
		PairedBool:             pairedBool,
		Connected:              connected,
		InitialEstimateTimeSec: 180,
		Networks:               GetConnectedNetwork(),
		SSPUrl:                 device.GetApiBaseUrl(),
		NewBindProcessSupport:  true,
	}

	for i, n := range results.Networks {
		logger.AppLogger().Debugf("ServiceInit, results.Networks[%v/%v]:%+v", i, len(results.Networks), n)
	}

	info := clientinfo.GetAdminPairedInfo()
	if info != nil {
		results.ClientUuid = info.ClientUuid
		results.BoxName = info.BoxName
	}

	results.GenerationEn = dtodevice.GetGenerationEn()
	//results.OpenSource = true
	results.DeviceName = dtodevice.GetDeviceName()
	results.DeviceNameEn = dtodevice.GetDeviceNameEn()

	results.GenerationEn = dtodevice.GetGenerationEn()
	results.GenerationEn = dtodevice.GetGenerationZh()

	results.ProductModel = dtodevice.GetProductModel()

	results.DeviceAbility = device_ability.GetAbilityModel()

	if device_ability.GetAbilityModel().RunInDocker {
		results.InitialEstimateTimeSec = int(config.Config.Box.InitialEstimateTimeSecRunInDocker)
	}

	return encwrapper.Enc(results)
}
