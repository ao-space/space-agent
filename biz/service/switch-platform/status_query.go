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
 * @Author: zhongguang
 * @Date: 2022-11-23 09:05:41
 * @Last Modified by: zhongguang
 * @Last Modified time: 2022-11-24 17:29:26
 */

package switchplatform

import (
	"agent/biz/model/dto"
	modelsp "agent/biz/model/switch-platform"
	"agent/biz/service/encwrapper"
	"errors"
	"fmt"

	"agent/utils/logger"
)

// 切换盒子对接的空间平台
func ServiceSwitchStatusQuery(req *modelsp.SwitchStatusQueryReq) (dto.BaseRspStr, error) {
	// logger.AppLogger().Debugf("ServiceSwitchStatusQuery, req:%+v", req)

	err := encwrapper.Check()
	if err != nil {
		logger.AppLogger().Warnf("check, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}

	rt, err := encwrapper.Dec(req.TransId)
	if err != nil {
		logger.AppLogger().Warnf("dec, err:%+v", err)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, nil
	}
	TransId := rt[0]

	logger.AppLogger().Debugf("transId:%+v", TransId)
	rsp, err := doStatusQuery(TransId)

	logger.AppLogger().Debugf("ServiceSwitchStatusQuery, transId:%v, rsp:%+v, err:%v", TransId, rsp, err)
	return rsp, err
}

func doStatusQuery(transId string) (dto.BaseRspStr, error) {
	logger.AppLogger().Debugf("doStatusQuery,transId=%v",
		transId)

	mtx.Lock()
	defer mtx.Unlock()
	if si == nil || si.TransId != transId {

		var basersp dto.BaseRspStr
		basersp.Code = dto.AgentCodeSwitchTaskNotFoundErr
		basersp.Message = fmt.Sprintf("transId(%v) not found.", transId)
		return basersp, errors.New(basersp.Message)

	}

	var resp modelsp.SwitchStatusQueryResp
	resp.Status = si.Status
	resp.TransId = si.TransId
	resp.UserDomain, _ = si.ImigrateResult.GetAdminDomain()

	return encwrapper.Enc(&resp)

}
