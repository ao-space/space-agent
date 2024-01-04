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

package pair

import (
	"agent/biz/docker"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	"agent/biz/model/dto/pair/tryout"
	"agent/config"
	"fmt"
	"net/http"
	"time"

	"agent/utils/logger"

	utilshttp "agent/utils/network/http"

	"github.com/dungeonsnd/gocom/encrypt/random"
)

func ServiceTryout(req *tryout.TryoutCodeReq) (dto.BaseRspStr, error) {
	if device_ability.GetAbilityModel().RunInDocker {
		if docker.ContainersDownloading == docker.GetDockerStatus() {
			err := fmt.Errorf("docker images is downloading...")
			logger.AppLogger().Warnf("ServiceTryout,%v", err)
			return dto.BaseRspStr{Code: dto.AgentCodeDockerPulling, Message: err.Error(), Results: nil}, err
		}
	}

	return presetBoxInfo(req)
}

// 预置试用信息
func presetBoxInfo(req *tryout.TryoutCodeReq) (dto.BaseRspStr, error) {

	// 平台请求结构
	type platformReqStruct struct {
		Email   string `json:"email"`
		Code    string `json:"code"`
		Type    string `json:"type"`
		BoxInfo struct {
			BoxUUID   string            `json:"boxUUID"`
			Desc      string            `json:"desc"`
			Extra     map[string]string `json:"extra"`
			BoxPubKey string            `json:"boxPubKey"`
			AuthType  string            `json:"authType"`
		} `json:"boxInfo"`
	}
	// 平台响应结构
	type platformRspStruct struct {
		Code      string `json:"code"`
		Message   string `json:"message"`
		RequestId string `json:"requestId"`

		State   int32 `json:"state"` // 0-正常;1-禁用;2-已过期
		BoxInfo struct {
			AuthType     string `json:"authType"`
			SnNumber     string `json:"snNumber"`
			IsRegistered bool   `json:"isRegistered"`
		} `json:"boxInfo"`
	}

	// 请求平台
	parms := &platformReqStruct{}
	parms.Email = req.Email
	parms.Code = req.TryoutCode
	parms.Type = "pc_open"
	parms.BoxInfo.BoxUUID = device.GetDeviceInfo().BoxUuid
	parms.BoxInfo.Desc = "pc tryout"
	parms.BoxInfo.Extra = make(map[string]string)
	// parms.BoxInfo.BoxPubKey = strings.ReplaceAll(string(device.GetBoxPubKey()), "\n", "")
	parms.BoxInfo.BoxPubKey = string(device.GetDevicePubKey())
	parms.BoxInfo.AuthType = "box_pub_key"
	url := device.GetApiBaseUrl() + config.Config.Platform.PresetBoxInfo.Path
	logger.AppLogger().Debugf("presetBoxInfo, url:%+v, parms:%+v", url, parms)

	var headers = map[string]string{"Request-Id": random.GenUUID()}
	var rsp platformRspStruct

	tryTotal := 6
	// var httpReq *http.Request
	var httpRsp *http.Response
	var body []byte
	var err1 error
	for i := 0; i < tryTotal; i++ {
		_, httpRsp, body, err1 = utilshttp.PostJsonWithHeaders(url, parms, headers, &rsp)
		if err1 != nil {
			// logger.AppLogger().Warnf("Failed PostJson, err:%v, @@httpReq:%+v, @@httpRsp:%+v, @@body:%v", err1, httpReq, httpRsp, string(body))
			if i == tryTotal-1 {
				return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err1.Error(), Results: nil}, err1
			}
			time.Sleep(time.Second * 2)
			continue
		} else {
			break
		}
	}

	logger.AppLogger().Infof("presetBoxInfo, rsp:%+v", rsp)
	logger.AppLogger().Infof("presetBoxInfo, httpRsp:%+v", httpRsp)
	logger.AppLogger().Infof("presetBoxInfo, body:%v", string(body))

	if httpRsp.StatusCode == http.StatusOK {
		if rsp.State == 1 {
			err1 := fmt.Errorf("httpRsp.StatusCode: %+v, rsp.State: %+v", httpRsp.StatusCode, rsp.State)
			rsp := dto.BaseRspStr{Code: dto.AgentCodeTryOutCodeDisabled, Message: err1.Error(), Results: nil}
			return rsp, nil
		} else if rsp.State == 2 {
			err1 := fmt.Errorf("httpRsp.StatusCode: %+v, rsp.State: %+v", httpRsp.StatusCode, rsp.State)
			rsp := dto.BaseRspStr{Code: dto.AgentCodeTryOutCodeExpired, Message: err1.Error(), Results: nil}
			return rsp, nil
		}

		err := device.UpdateSnNumber(rsp.BoxInfo.SnNumber)
		if err != nil {
			err1 := fmt.Errorf("failed UpdateSnNumber: %+v", err)
			rsp := dto.BaseRspStr{Code: dto.AgentCodeTryOutCodeExpired, Message: err1.Error(), Results: nil}
			return rsp, nil
		}

		device.UpdateApplyEmail(req.Email)

		rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr, Message: "OK", Results: nil}
		return rsp, nil

	} else if httpRsp.StatusCode == http.StatusBadRequest {
		c := dto.AgentCodeBadReqStr
		if rsp.Code == "PSP-2047" {
			c = dto.AgentCodeTryOutCodeError
		} else if rsp.Code == "PSP-2052" {
			c = dto.AgentCodeTryOutCodeHasUsed
		}
		err1 := fmt.Errorf("httpRsp.StatusCode: %+v, rsp.Code: %+v, rsp.Message: %+v", httpRsp.StatusCode, rsp.Code, rsp.Message)
		rsp := dto.BaseRspStr{Code: c, Message: err1.Error(), Results: nil}
		return rsp, err1

	} else {
		err1 := fmt.Errorf("httpRsp.StatusCode: %+v, rsp.Code: %+v, rsp.Message: %+v", httpRsp.StatusCode, rsp.Code, rsp.Message)
		rsp := dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr, Message: err1.Error(), Results: nil}
		return rsp, err1
	}
}
