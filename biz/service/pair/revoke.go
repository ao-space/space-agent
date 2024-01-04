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
 * @Date: 2021-12-14 15:11:52
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-25 14:02:03
 * @Description:
 */

package pair

import (
	"agent/biz/model/dto"
	dtopair "agent/biz/model/dto/pair"
	"agent/biz/service/call"
	"agent/biz/service/encwrapper"
	"agent/config"

	"agent/utils/logger"
)

// 管理员解绑
func ServiceRevoke(req *dtopair.RevokeReq) (dto.BaseRspStr, error) {

	err := encwrapper.Check()
	if err != nil {
		logger.AppLogger().Warnf("check, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}

	rt, err := encwrapper.Dec(req.Password)
	if err != nil {
		logger.AppLogger().Warnf("dec, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}
	password := rt[0]
	clientUUID := ""

	if len(req.ClientUUID) > 0 {
		rt, err := encwrapper.Dec(req.ClientUUID)
		if err != nil {
			logger.AppLogger().Warnf("dec, err:%+v", err)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err.Error()}, nil
		}
		clientUUID = rt[0]
	}

	return doRevoke(password, clientUUID)
}

func doRevoke(password string, clientUUID string) (dto.BaseRspStr, error) {
	var results call.MicroServerRsp
	reqMap := make(map[string]string)
	reqMap["passcode"] = password
	reqMap["userid"] = "1"
	if len(clientUUID) > 0 {
		reqMap["clientUUID"] = clientUUID
	}

	resp, err := call.CallServiceByForm("POST", config.Config.Account.AdminRevoke.Url, reqMap, &results)
	if err != nil {
		logger.AppLogger().Warnf("failed CallServiceByForm, err:%+v, resp:%+v", err, resp)
		return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err.Error()},
			err
	}
	return encwrapper.Enc(results)
}
