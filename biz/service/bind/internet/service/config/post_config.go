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
	"agent/biz/docker"
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/dto"
	"agent/biz/model/dto/bind/internet/service/config"
	"agent/biz/model/gt"
	"agent/biz/service/base"
	"agent/biz/service/call"
	"agent/biz/service/pair"
	agentconfig "agent/config"
	"agent/utils"
	"agent/utils/logger"
	"agent/utils/retry"
	"fmt"
	"time"
)

type InternetServiceConfig struct {
	base.BaseService
	PairedInfo *clientinfo.AdminPairedInfo
	GTConfig   *gt.Config
}

func (svc *InternetServiceConfig) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("InternetServiceConfig Process, svc.RequestId:%v", svc.RequestId)

	svc.PairedInfo = clientinfo.GetAdminPairedInfo()
	//paired := clientinfo.GetAdminPairedStatus()
	if !svc.PairedInfo.AlreadyBound() {
		err := fmt.Errorf("check, paired:%+v", svc.PairedInfo.Status())
		return dto.BaseRspStr{Code: dto.AgentCodeUnpairedBeforeStr, Message: err.Error()}
	}

	if svc.Req == nil {
		err := fmt.Errorf("req is nil")
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, RequestId: svc.RequestId, Message: err.Error()}
	}

	req := svc.Req.(*config.ConfigReq)
	// logger.AppLogger().Debugf("InternetServiceConfig Process, req:%+v", req)
	if device.GetConfig().EnableInternetAccess { // 之前处于开启状态, 准备关闭
		if req.EnableInternetAccess {
			return svc.BaseService.Process()
		}

		err := docker.StopContainerImmediately(agentconfig.Config.Docker.NetworkClientContainerName)
		if err != nil {
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, RequestId: svc.RequestId, Message: err.Error()}
		}

	} else { // 之前处于关闭状态, 准备开启
		if !req.EnableInternetAccess {
			return svc.BaseService.Process()
		}
		if req.PlatformApiBase != "" {
			device.SetApiBaseUrl(req.PlatformApiBase)
		}
		// 调用网关接口切换平台
		var microServerRsp call.MicroServerRsp
		if req.PlatformApiBase != agentconfig.Config.Platform.APIBase.Url {
			type ChangePlatformReq struct {
				SSPlatformUrl string `json:"ssplatformUrl"`
			}
			changePlatformReq := &ChangePlatformReq{SSPlatformUrl: req.PlatformApiBase}
			switchPlatformReq := func() error {
				err := call.CallServiceByPost(agentconfig.Config.GateWay.SwitchPlatform.Url, nil, changePlatformReq, &microServerRsp)
				if err != nil {
					return err
				}
				return nil
			}
			err := retry.Retry(switchPlatformReq, 3, time.Second*2)
			if err != nil {
				logger.AppLogger().Errorf("switch platform error:%v", err)
				return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err.Error()}
			}
		}
		// 1. 调用平台注册盒子接口
		err := pair.ServiceRegisterBox()
		if err != nil {
			host, _, _ := utils.ParseUrl(device.GetApiBaseUrl())
			hostOfficial, _, _ := utils.ParseUrl(agentconfig.Config.Platform.APIBase.Url)
			if host != hostOfficial {
				err1 := fmt.Errorf("ServiceRegisterBox failed, current space platform:%v, %+v", host, err)
				logger.AppLogger().Warnf("%v", err1)
				return dto.BaseRspStr{Code: dto.AgentCodePrivateSSPRegBoxErr, Message: err1.Error()}
			} else {
				err1 := fmt.Errorf("ServiceRegisterBox failed, %+v", err)
				logger.AppLogger().Warnf("%v", err1)
				return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err1.Error()}
			}
		}
		device.SetBoxRegistered()
		if err != nil {
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, RequestId: svc.RequestId, Message: err.Error()}
		}
		svc.GTConfig = &gt.Config{}
		err = svc.GTConfig.Init()
		if err != nil {
			logger.AppLogger().Errorf("failed to init gt config yaml:%v", err)
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, RequestId: svc.RequestId, Message: err.Error()}
		}
		err = docker.DockerUpImmediately(nil)
		if err != nil {
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, RequestId: svc.RequestId, Message: err.Error()}
		}
	}

	device.SetConfig(&device.InternetServiceConfig{EnableInternetAccess: req.EnableInternetAccess})

	// 调用网关 开启互联网服务的接口

	type ChannelInfoStruct struct {
		Wan bool `json:"wan,omitempty"`
	}
	microServerReq := &ChannelInfoStruct{Wan: req.EnableInternetAccess}
	queryParms := fmt.Sprintf("?userId=%v&AccessToken-clientUUID=%v", "1", req.ClientUUID)
	var microServerRsp call.MicroServerRsp
	internetAccessReq := func() error {
		err := call.CallServiceByPost(agentconfig.Config.Account.NetworkChannelWan.Url+queryParms, nil, microServerReq, &microServerRsp)
		if err != nil {
			return err
		}
		return nil
	}
	err := retry.Retry(internetAccessReq, 3, time.Second*2)
	if err != nil {
		return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err.Error()}
	}
	logger.AppLogger().Debugf("microServerRsp:%+v", microServerRsp)
	if microServerRsp.Code != dto.GatewayCodeOkStr && microServerRsp.Code != dto.AccountCodeOkStr {
		return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr,
			Message: fmt.Errorf("call microsever failed, code:%v, message:%v", microServerRsp.Code, microServerRsp.Message).Error()}
	}

	rsp := &config.ConfigRsp{EnableInternetAccess: device.GetConfig().EnableInternetAccess,
		EnableP2P:  device.GetConfig().EnableInternetAccess,
		EnableLAN:  true,
		UserDomain: svc.PairedInfo.AdminDomain()}
	rsp.ConnectedNetwork = pair.GetConnectedNetwork()
	svc.Rsp = rsp
	return svc.BaseService.Process()
}
