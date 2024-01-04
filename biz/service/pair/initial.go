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
 * @Date: 2021-12-03 21:50:56
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-25 14:01:41
 * @Description:
 */
package pair

import (
	"agent/biz/model/dto"
	dtopair "agent/biz/model/dto/pair"
	"agent/biz/service/call"
	"agent/biz/service/encwrapper"
	"agent/config"
	"fmt"
	"time"

	"agent/utils/logger"
)

func ServiceInitial(req *dtopair.PasswordInfo) (dto.BaseRspStr, error) {
	logger.AppLogger().Debugf("ServiceInitial")
	// logger.AccessLogger().Debugf("[ServiceInitial], req:%+v", req)
	err := encwrapper.Check()
	if err != nil {
		err1 := fmt.Errorf("check failed, err:%v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err1.Error()}, err1
	}

	rt, err := encwrapper.Dec(req.Password)
	if err != nil {
		logger.AppLogger().Warnf("dec, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}
	password := rt[0]

	// 等待 account 的 initial 接口调用成功
	var err1 error
	var rsp dto.BaseRspStr
	for i := 0; i < 3; i++ {
		time.Sleep(time.Duration(3) * time.Second)
		rsp, err1 = doInitial(password)
		if err1 == nil {
			break
		}
		logger.AppLogger().Debugf("ServiceInitial, waiting docker account service initial return, err1:%v", err1)
	}
	logger.AppLogger().Debugf("ServiceInitial, loop break, err1:%v", err1)
	// account 的 initial 接口一定时间内没调用成功
	if err1 != nil {
		logger.AppLogger().Debugf("ServiceInitial, return rsp:%+v, err1:%v", rsp, err1)
		return rsp, err1
	}

	// 成功返回
	logger.AppLogger().Debugf("ServiceInitial return rsp:%+v, rsp.Results:%+v", rsp, rsp.Results)
	return rsp, nil
}

func doInitial(password string) (dto.BaseRspStr, error) {
	var results call.MicroServerRsp
	reqString := fmt.Sprintf("flag=true&password=%v", password)

	resp, err := call.CallServiceByFormStr("POST", config.Config.Account.AdminInitial.Url, reqString, &results)
	if err != nil {
		logger.AppLogger().Warnf("failed CallServiceByFormBool, err:%+v, resp:%+v", err, resp)
		return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err.Error()},
			err
	}
	if results.Code != "ACC-200" {
		err1 := fmt.Errorf("CallServiceByFormBool, results.Code:%+v", results.Code)
		logger.AppLogger().Warnf("%v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeCallServiceFailedStr, Message: err1.Error()},
			err1
	}

	return encwrapper.Enc(results)
}
