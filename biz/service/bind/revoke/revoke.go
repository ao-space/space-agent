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

package revoke

import (
	"agent/biz/model/dto"
	"agent/biz/model/dto/bind/revoke"
	"agent/biz/service/base"
	"agent/biz/service/call"
	"agent/config"
	"fmt"

	"agent/utils/logger"
)

type RevokeService struct {
	base.BaseService
}

func (svc *RevokeService) Process() dto.BaseRspStr {
	logger.AppLogger().Debugf("RevokeService Process")

	// paired := clientinfo.GetAdminPairedStatus()
	// if paired != clientinfo.ClientPairedStatusBind {
	// 	err := fmt.Errorf("check, paired:%+v", paired)
	// 	return dto.BaseRspStr{Code: dto.AgentCodeUnpairedBeforeStr, Message: err.Error()}
	// }

	if svc.Req == nil {
		err := fmt.Errorf("req is nil")
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr, RequestId: svc.RequestId, Message: err.Error()}
	}

	req := svc.Req.(*revoke.RevokeReq)
	// logger.AppLogger().Debugf("RevokeService Process, req:%+v", req)
	microServerRsp, err := doRevoke(req.Password, req.ClientUuid)
	if err != nil {
		logger.AppLogger().Debugf("%v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err.Error()}
	}

	rsp := &revoke.RevokeRsp{Code: microServerRsp.Code,
		RequestId: microServerRsp.RequestId,
		Message:   microServerRsp.Message,
		Results:   microServerRsp.Results}

	if microServerRsp.Code == dto.GatewayCodeOkStr || microServerRsp.Code == dto.AccountCodeOkStr {
		jwtToken, err := base.CreateAgentToken(req.ClientUuid)
		if err != nil {
			logger.AppLogger().Debugf("%v", err)
			return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err.Error()}
		}
		rsp.AgentToken = jwtToken
	}

	svc.Rsp = rsp

	return svc.BaseService.Process()
}

func doRevoke(password string, clientUUID string) (*call.MicroServerRsp, error) {
	var results call.MicroServerRsp
	reqMap := make(map[string]string)
	reqMap["passcode"] = password
	reqMap["userid"] = "1"
	reqMap["clientUUID"] = clientUUID

	resp, err := call.CallServiceByForm("POST", config.Config.Account.AdminRevoke.Url, reqMap, &results)
	if err != nil {
		logger.AppLogger().Warnf("failed CallServiceByForm, err:%+v, resp:%+v", err, resp)
		return nil, err
	}
	return &results, nil
}
