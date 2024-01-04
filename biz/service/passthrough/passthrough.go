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

package passthrough

import (
	agentrouters "agent/biz/model/agent-routers"
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/passthrough"
	"agent/biz/service/base"
	"agent/biz/service/encwrapper"
	"agent/config"
	"fmt"

	"agent/utils/logger"

	utilshttp "agent/utils/network/http"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"github.com/dungeonsnd/gocom/encrypt/random"
)

type PassthroughService struct {
	base.BaseService
}

func (svc *PassthroughService) Process() dto.BaseRspStr {
	req := svc.Req.(*passthrough.PassthroughReq)

	if len(req.ApiPath) < 1 {
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: fmt.Sprintf("req.ApiPath length error")}
	}

	// set headers
	newHeaders := make(map[string]string)
	for k, v := range svc.Header {
		newHeaders[k] = v[0]
	}

	// set clientUuid
	clientUUID := ""
	pairedInfo := clientinfo.GetAdminPairedInfo()
	if pairedInfo == nil {
		return dto.BaseRspStr{Code: dto.AgentCodeUnpairedBeforeStr,
			Message: fmt.Sprintf("not paired, pairedInfo=%+v", pairedInfo)}
	} else {
		clientUUID = pairedInfo.ClientUuid
	}
	if len(clientUUID) > 0 {
		newHeaders["clientUuid"] = clientUUID
	}

	// 已经解绑时，account 容器会删除 /etc/meta/admin/admin 中的 clientUuid，导致取不到之前绑定端的 clientUuid。这里再传一下 btid。
	btid := device.GetDeviceInfo().Btid
	newHeaders["btid"] = btid

	//generate sign
	if device_ability.GetAbilityModel().SecurityChipSupport {
		if len(clientUUID) > 0 {
			clientUuidSign, err := device.SignFromSecurityChip([]byte(clientUUID))
			if err != nil {
				logger.AppLogger().Warnf("%+v", clientUuidSign)
				return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
					Message: err.Error()}
			}
			newHeaders["clientUuidSign"] = clientUuidSign
		}

		b64SignedBtid, err := device.SignFromSecurityChip([]byte(btid))
		if err != nil {
			logger.AppLogger().Warnf("%+v", b64SignedBtid)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err.Error()}
		}
		newHeaders["btidSign"] = b64SignedBtid

	} else {
		pri, err := encwrapper.GetPrivateKey(string(device.GetDevicePriKey()))
		if err != nil {
			err1 := fmt.Errorf("%+v", err)
			logger.AppLogger().Warnf("%+v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err1.Error()}
		}

		if len(clientUUID) > 0 {
			d, err := encwrapper.Sign(pri, []byte(clientUUID))
			if err != nil {
				err1 := fmt.Errorf("%+v", err)
				logger.AppLogger().Warnf("%+v", err1)
				return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
					Message: err1.Error()}
			}
			newHeaders["clientUuidSign"] = encoding.Base64Encode(d)
		}

		d, err := encwrapper.Sign(pri, []byte(btid))
		if err != nil {
			err1 := fmt.Errorf("%+v", err)
			logger.AppLogger().Warnf("%+v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err1.Error()}
		}
		newHeaders["btidSign"] = encoding.Base64Encode(d)
	}

	newHeaders["Content-Type"] = "application/json"
	if _, ok := newHeaders["Request-Id"]; !ok {
		newHeaders["Request-Id"] = random.GenUUID()
	}

	if agentrouters.AgentRouters == nil || agentrouters.AgentRouters.ValidPaths == nil {
		err1 := fmt.Errorf("AgentRouters error")
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err1.Error()}
	}

	validPaths, found := agentrouters.AgentRouters.ValidPaths[req.ServiceName]

	if found {
		reqPath := "/" + req.ApiVersion + req.ApiPath

		currentServiceValidPaths := make(map[string]string)
		for _, v := range validPaths {
			currentServiceValidPaths[v] = ""
		}
		_, validReq := currentServiceValidPaths[reqPath]
		if !validReq {
			err1 := fmt.Errorf("illegal reqPath:%v. ServiceName:%v",
				reqPath, req.ServiceName)
			logger.AppLogger().Warnf("%+v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err1.Error()}
		}

		url := config.Config.GateWay.APIRoot.Url + reqPath
		logger.AppLogger().Debugf("---- ServicePassthrough, url:%+v", url)
		var callRsp interface{}
		_, httpRsp, rspBody, err1 := utilshttp.PostJsonWithHeaders(url,
			req.Entity, newHeaders, callRsp)
		if err1 != nil {
			// logger.AppLogger().Warnf("Failed CallServiceByPost, err:%v, @@httpReq:%+v, @@httpRsp:%+v, @@rspBody:%v", err1, httpReq, httpRsp, string(rspBody))
			return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err1.Error()}
		} else {
			logger.AppLogger().Debugf("ServicePassthrough, httpRsp:%+v", httpRsp)
			// logger.AppLogger().Debugf("ServicePassthrough, callRsp:%+v", callRsp)
		}

		svc.RequestId = newHeaders["Request-Id"]

		svc.RspBytes = rspBody
		return svc.BaseService.Process()

	} else {
		err1 := fmt.Errorf("unsupported ServiceName: %+v", req.ServiceName)
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err1.Error()}
	}
}
