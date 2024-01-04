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

package create

import (
	"agent/biz/docker"
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/did"
	"agent/biz/model/did/leveldb"
	"agent/biz/model/dto"
	"agent/biz/model/dto/bind/space/create"
	"agent/biz/model/gt"
	"agent/biz/service/base"
	"agent/biz/service/call"
	"agent/biz/service/pair"
	"agent/config"
	"agent/utils"
	"agent/utils/logger"
	"agent/utils/retry"
	"encoding/base64"
	"fmt"
	"time"
)

type SpaceCreateService struct {
	base.BaseService
	PairedInfo *clientinfo.AdminPairedInfo
	GTConfig   *gt.Config
}

func (svc *SpaceCreateService) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("SpaceCreateService Process, svc.RequestId:%v", svc.RequestId)

	var err error
	svc.PairedInfo = clientinfo.GetAdminPairedInfo()
	if svc.PairedInfo.Status() == clientinfo.DeviceAlreadyBound {
		err := fmt.Errorf("check, paired:%+v", svc.PairedInfo.Status())
		return dto.BaseRspStr{Code: dto.AgentCodeAlreadyPairedStr, Message: err.Error()}
	}

	if svc.Req == nil {
		err := fmt.Errorf("req is nil")
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, RequestId: svc.RequestId, Message: err.Error()}
	}
	rsp := &create.CreateRsp{}

	req := svc.Req.(*create.CreateReq)
	// logger.AppLogger().Debugf("SpaceCreateService Process, req:%+v", req)
	logger.AppLogger().Debugf("SpaceCreateService Process, pairedInfo:%+v", svc.PairedInfo)
	logger.AppLogger().Debugf("SpaceCreateService, rebind:%+v", svc.PairedInfo.Rebind())

	// 注册平台
	if req.EnableInternetAccess && !svc.PairedInfo.Rebind() { // 首次且开启互联网通道情况下。解绑后重新绑定时不改变互联网服务配置。

		if req.PlatformApiBase != "" {
			device.SetApiBaseUrl(req.PlatformApiBase)
			//envFiles := make(map[string]map[string]string)
			//gwEnvFile := make(map[string]string)
			//gwEnvFile["APP_SSPLATFORM_URL"] = req.PlatformApiBase
			//envFiles["aospace-gateway.env"] = gwEnvFile
			//docker.ProcessEnv(config.Config.Docker.ComposeFile, envFiles)
		}

		result, err := svc.registerDevice(req)
		if err != nil {
			return *result
		}
		logger.AppLogger().Debugf("registerDevice, result:%+v", result)
	}
	if !svc.PairedInfo.Rebind() { // 首次
		logger.AppLogger().Debugf("SpaceCreateService, SetConfig req.EnableInternetAccess:%+v", req.EnableInternetAccess)
		device.SetConfig(&device.InternetServiceConfig{EnableInternetAccess: req.EnableInternetAccess})
	}

	rsp.ConnectedNetwork = pair.GetConnectedNetwork()

	logger.AppLogger().Debugf("SpaceCreateService, BeginTransaction")
	// 创建 did
	var levelDBTrans *leveldb.Trans
	if req.VerifyMethod != nil {
		levelDBTrans, err = leveldb.BeginTransaction() // 开启事务
		defer levelDBTrans.Rollback()                  // 退出时回滚事务. 如果成功, 函数返回之前主动 commit.
		didDocBytes, encryptedPriKeyBytes, err := svc.createDid(levelDBTrans, req)
		logger.AppLogger().Debugf("createDid, len(didDocBytes):%v, len(encryptedPriKeyBytes):%v, err: %v",
			len(didDocBytes), len(encryptedPriKeyBytes), err)
		if err != nil {
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
		}

		rsp.DIDDoc = base64.StdEncoding.EncodeToString(didDocBytes)
		rsp.EncryptedPriKeyBytes = base64.StdEncoding.EncodeToString(encryptedPriKeyBytes)
	}

	// 创建 agent token
	jwtToken, err := base.CreateAgentToken(req.ClientUuid)
	if err != nil {
		logger.AppLogger().Debugf("%v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
	}
	logger.AppLogger().Debugf("CreateAgentToken, jwtToken:%v", jwtToken)

	clientinfo.SaveClientExchangeKey()

	// 调用网关接口创建账号
	microServerRsp, err := svc.callGateway(req)
	if err != nil {
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
	}
	rsp.AgentToken = jwtToken
	rsp.EnableInternetAccess = device.GetConfig().EnableInternetAccess
	rsp.SpaceUserInfo = microServerRsp
	logger.AppLogger().Debugf("callGateway, microServerRsp:%v", microServerRsp)

	// TODO: update did doc history
	// ...
	// 启动/停止转发服务
	svc.startOrStopContainers()
	logger.AppLogger().Debugf("startOrStopContainers")

	svc.Rsp = rsp

	if levelDBTrans != nil {
		levelDBTrans.Commit() // commit 提交事务
		logger.AppLogger().Debugf("callGateway, Commit")
	}
	return svc.BaseService.Process()
}

func (svc *SpaceCreateService) registerDevice(req *create.CreateReq) (*dto.BaseRspStr, error) {
	logger.AppLogger().Debugf("registerDevice, req:%+v ", req)

	// 调用平台注册盒子接口
	err := pair.ServiceRegisterBox()
	if err != nil {
		host, _, _ := utils.ParseUrl(device.GetApiBaseUrl())
		hostOfficial, _, _ := utils.ParseUrl(config.Config.Platform.APIBase.Url)
		if host != hostOfficial {
			err1 := fmt.Errorf("ServiceRegisterBox failed, current space platform:%v, %+v", host, err)
			logger.AppLogger().Warnf("%v", err1)
			return &dto.BaseRspStr{Code: dto.AgentCodePrivateSSPRegBoxErr, Message: err1.Error()}, err1
		} else {
			err1 := fmt.Errorf("ServiceRegisterBox failed, %+v", err)
			logger.AppLogger().Warnf("%v", err1)
			return &dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err1.Error()}, err1
		}
	}
	device.SetBoxRegistered()

	return nil, nil
}

func (svc *SpaceCreateService) callGateway(req *create.CreateReq) (call.MicroServerRsp, error) {
	// logger.AppLogger().Debugf("callGateway, req:%+v ", req)

	var microServerRsp call.MicroServerRsp

	// 调用网关接口切换平台

	if req.PlatformApiBase != config.Config.Platform.APIBase.Url && req.PlatformApiBase != "" {
		type ChangePlatformReq struct {
			SSPlatformUrl string `json:"ssplatformUrl"`
		}
		changePlatformReq := &ChangePlatformReq{SSPlatformUrl: req.PlatformApiBase}
		switchPlatformReq := func() error {
			err := call.CallServiceByPost(config.Config.GateWay.SwitchPlatform.Url, nil, changePlatformReq, &microServerRsp)
			if err != nil {
				return err
			}
			return nil
		}
		err := retry.Retry(switchPlatformReq, 3, time.Second*2)
		if err != nil {
			return microServerRsp, err
		}
	}
	// 调用网关的 "/space/v2/api/space/admin"
	type CreateStruct struct {
		ClientUUID           string `json:"clientUUID,omitempty"`
		PhoneModel           string `json:"phoneModel,omitempty"`
		ApplyEmail           string `json:"applyEmail,omitempty"`
		SpaceName            string `json:"spaceName,omitempty"`
		Password             string `json:"password,omitempty"`
		EnableInternetAccess bool   `json:"enableInternetAccess"`
	}
	microServerReq := &CreateStruct{ClientUUID: req.ClientUuid,
		PhoneModel:           req.ClientPhoneModel,
		Password:             req.Password,
		EnableInternetAccess: req.EnableInternetAccess}
	if !svc.PairedInfo.Rebind() { // 解绑后重新绑定时不改变 SpaceName。
		microServerReq.SpaceName = req.SpaceName
	}

	applyEmail, _ := device.GetApplyEmail()
	if len(applyEmail) > 0 {
		microServerReq.ApplyEmail = applyEmail
	}

	createAdminReq := func() error {
		err := call.CallServiceByPost(config.Config.Account.SpaceAdmin.Url, nil, microServerReq, &microServerRsp)
		if err != nil {
			return err
		}
		return nil
	}
	err := retry.Retry(createAdminReq, 3, time.Second*2)
	if err != nil {
		return microServerRsp, err
	}

	logger.AppLogger().Debugf("microServerRsp:%+v", microServerRsp)
	if microServerRsp.Code != dto.GatewayCodeOkStr && microServerRsp.Code != dto.AccountCodeOkStr {
		return microServerRsp, fmt.Errorf("call microsever failed, code:%v, message:%v",
			microServerRsp.Code, microServerRsp.Message)
	}

	return microServerRsp, nil
}

func (svc *SpaceCreateService) createDid(levelDBTrans *leveldb.Trans, req *create.CreateReq) ([]byte, []byte, error) {
	logger.AppLogger().Debugf("createDid, req:%+v ", req)

	// aoId, err := getAoId(microServerRsp.Results)
	// if err != nil {
	// 	logger.AppLogger().Warnf(err.Error())
	// 	return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
	// }
	// logger.AppLogger().Debugf("aoId:%+v", aoId)

	aoId := "aoid-1"

	encryptedPriKeyBytes, didDocBytes, didStr, err := did.CreateDocument(levelDBTrans, aoId, req.Password, req.VerifyMethod)
	if err != nil {
		err1 := fmt.Errorf("failed CreateDocument, err: %v", err)
		logger.AppLogger().Warnf(err1.Error())
		return nil, nil, err1
	}
	logger.AppLogger().Debugf("CreateDocument")

	err = did.SaveAoIdToDid(levelDBTrans, aoId, didStr)
	if err != nil {
		err1 := fmt.Errorf("failed SaveAoIdToDidString, err: %v", err)
		logger.AppLogger().Warnf(err1.Error())
		return nil, nil, err1
	}
	logger.AppLogger().Debugf("SaveAoIdToDid succ, %+v->%v", aoId, didStr)

	return didDocBytes, encryptedPriKeyBytes, nil
}

func (svc *SpaceCreateService) startOrStopContainers() {
	logger.AppLogger().Debugf("startOrStopContainers")

	if device.GetConfig().EnableInternetAccess {
		svc.GTConfig = &gt.Config{}
		err := svc.GTConfig.Init()
		if err != nil {
			logger.AppLogger().Errorf("failed to init gt config yaml:%v", err)
		}
		err = docker.DockerUpImmediately(nil)
		if err != nil {
			logger.AppLogger().Warnf("DockerUpImmediately: %v", err)
		}
	} else {
		go func() {
			err := docker.StopContainerImmediately(config.Config.Docker.NetworkClientContainerName)
			if err != nil {
				logger.AppLogger().Warnf("StopContainerImmediately: %v", err)
			}
		}()
	}
}

func getAoId(results interface{}) (string, error) {
	logger.AppLogger().Debugf("getAoId, results:%+v ", results)

	if results == nil {
		err1 := fmt.Errorf("getAoId,results is nil")
		logger.AppLogger().Warnf(err1.Error())
		return "", err1
	}

	logger.AppLogger().Debugf("getAoId, results : %+v", results)
	m := results.(map[string]interface{})
	logger.AppLogger().Debugf("m: %+v", m)
	if v, ok := m["aoId"]; ok {
		logger.AppLogger().Debugf("v: %+v", v)
		aoId := v.(string)
		if len(aoId) < 1 {
			err1 := fmt.Errorf("getAoId, aoId:%v", aoId)
			logger.AppLogger().Warnf(err1.Error())
			return "", err1
		}
		return aoId, nil
	}

	return "", fmt.Errorf("aoId not found in %+v", m)
}
