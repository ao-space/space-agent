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
 * @Date: 2021-12-13 15:16:59
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-17 14:46:47
 * @Description:
 */

package pair

import (
	"agent/biz/model/dto"
	dtopair "agent/biz/model/dto/pair"
	"agent/biz/service/encwrapper"
	"fmt"

	"agent/utils/logger"
)

func ServiceWifiList(req *dtopair.WifiListReq, usingEncrypt bool) (dto.BaseRspStr, error) {

	if usingEncrypt {
		err := encwrapper.Check()
		if err != nil {
			err1 := fmt.Errorf("check failed, err:%v", err)
			logger.AppLogger().Warnf("%+v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err1.Error()}, err1
		}
	}

	results := GetWifiList()

	if usingEncrypt {
		return encwrapper.Enc(results)
	} else {
		return dto.BaseRspStr{Code: dto.AgentCodeOkStr,
			Message: "OK",
			Results: results}, nil
	}
}
