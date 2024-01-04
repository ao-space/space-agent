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
 * @Date: 2021-12-03 21:26:32
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-25 14:01:59
 * @Description:
 */
package pair

import (
	"agent/biz/model/dto"
	"agent/biz/service/call"
	"agent/biz/service/encwrapper"
	"agent/config"

	dtopair "agent/biz/model/dto/pair"

	"agent/utils/logger"
)

func ServiceSetPassword(req *dtopair.PasswordInfo) (dto.BaseRspStr, error) {

	err := encwrapper.Check()
	if err != nil {
		logger.AppLogger().Warnf("check, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}

	rt, err := encwrapper.Dec(req.Password)
	if err != nil {
		// logger.AppLogger().Warnf("dec, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}
	password := rt[0]

	return doSetPassword(password)
}

func doSetPassword(password string) (dto.BaseRspStr, error) {
	type CreateStruct struct {
		Password string `json:"password,omitempty"`
	}
	req := &CreateStruct{Password: password}
	var results call.MicroServerRsp
	err := call.CallServiceByPost(config.Config.Account.AdminSetPassword.Url, nil, req, &results)
	if err != nil {
		return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err.Error()},
			err
	}
	return encwrapper.Enc(results)
}
