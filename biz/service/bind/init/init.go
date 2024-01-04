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

package init

import (
	"agent/biz/alivechecker/model"
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/dto/bind/bindinit"
	dtodevice "agent/biz/model/dto/device"
	dtopair "agent/biz/model/dto/pair"
	"agent/biz/service/base"
	"agent/biz/service/pair"
	"agent/config"
	"agent/utils/logger"
	"agent/utils/version"
)

type InitService struct {
	base.BaseService
	PairedInfo *clientinfo.AdminPairedInfo
}

func (svc *InitService) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("InitService Process")
	req := svc.Req.(*bindinit.InitReq)
	// logger.AppLogger().Debugf("InitService Process, req:%+v", req)
	if req != nil && len(req.ClientUuid) > 0 && len(req.ClientVersion) > 0 {
		clientinfo.SetClientVersion(req.ClientUuid, req.ClientVersion)
	}

	svc.PairedInfo = clientinfo.GetAdminPairedInfo()
	// https://code.eulix.xyz/bp/cicada/-/issues/465#note_82345
	boxUuid := device.GetDeviceInfo().BoxUuid

	//paired := clientinfo.GetAdminPairedStatus()
	//pairedBool := clientinfo.GetAdminPairedStatus() == clientinfo.ClientPairedStatusBind
	//logger.AppLogger().Debugf("ServiceInit, paired:%v")

	connected := 1
	if model.Get().PingCloudHost {
		logger.AppLogger().Debugf("ServiceInit, connected")
		connected = 0
	} else {
		logger.AppLogger().Debugf("ServiceInit, not connected")
	}

	rsp := &dtopair.InitResult{
		BoxName:                "傲空间",
		ProductId:              boxUuid,
		BoxUuid:                boxUuid,
		Paired:                 svc.PairedInfo.Status(),
		PairedBool:             svc.PairedInfo.AlreadyBound(),
		Connected:              connected,
		InitialEstimateTimeSec: 180,
		Networks:               pair.GetConnectedNetwork(),
		SSPUrl:                 device.GetApiBaseUrl(),
		NewBindProcessSupport:  true,
	}

	for i, n := range rsp.Networks {
		logger.AppLogger().Debugf("ServiceInit, rsp.Networks[%v/%v]:%+v", i, len(rsp.Networks), n)
	}

	if svc.PairedInfo != nil {
		rsp.ClientUuid = svc.PairedInfo.ClientUuid
		rsp.BoxName = svc.PairedInfo.BoxName
	}
	boxVersion := version.GetInstalledAgentVersionRemovedNewLine()
	if len(boxVersion) < 3 {
		boxVersion = config.VersionNumber
	}
	rsp.SpaceVersion = boxVersion
	//rsp.OpenSource = true
	rsp.GenerationEn = dtodevice.GetGenerationEn()
	rsp.DeviceName = dtodevice.GetDeviceName()
	rsp.DeviceNameEn = dtodevice.GetDeviceNameEn()
	rsp.GenerationEn = dtodevice.GetGenerationEn()
	rsp.GenerationEn = dtodevice.GetGenerationZh()
	rsp.ProductModel = dtodevice.GetProductModel()
	rsp.DeviceAbility = device_ability.GetAbilityModel()
	if device_ability.GetAbilityModel().RunInDocker {
		rsp.InitialEstimateTimeSec = int(config.Config.Box.InitialEstimateTimeSecRunInDocker)
	}

	svc.Rsp = rsp
	return svc.BaseService.Process()
}
