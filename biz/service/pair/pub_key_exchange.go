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
 * @Date: 2021-12-13 13:53:01
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-24 13:53:17
 * @Description:
 */

package pair

import (
	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/model/dto"
	dtopair "agent/biz/model/dto/pair"
	"agent/biz/service/encwrapper"
	"fmt"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
)

func ServicePubKeyExchange(req *dtopair.PubKeyExchangeReq) (dto.BaseRspStr, error) {
	logger.AppLogger().Debugf("ServicePubKeyExchange, len(req.ClientPubKey):%+v", len(req.ClientPubKey))
	logger.AppLogger().Debugf("ServicePubKeyExchange, req.ClientPubKey:%+v", req.ClientPubKey)
	logger.AppLogger().Debugf("ServicePubKeyExchange, req.SignedBtid:%+v", req.SignedBtid)

	clientinfo.SetClientExchangePubKey(req.ClientPubKey)
	clientinfo.SetClientExchangePriKey(req.ClientPriKey)

	enableSigned := true
	b64SignedBtid := ""
	if enableSigned {
		// check sign
		pub, err := encwrapper.GetPublicKey(req.ClientPubKey)
		if err != nil {
			err1 := fmt.Errorf("%+v", err)
			logger.AppLogger().Warnf("%+v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err1.Error()}, err1
		}
		bSignedBtid, err := encoding.Base64Decode(req.SignedBtid)
		if err != nil {
			err1 := fmt.Errorf("Base64Decode(req.SignedBtid), %+v", err)
			logger.AppLogger().Warnf("%+v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err1.Error()}, err1
		}
		btid := device.GetDeviceInfo().Btid
		err = encwrapper.Verify(pub, bSignedBtid, []byte(btid))
		if err != nil {
			err1 := fmt.Errorf("verify NOT PASSED! %+v", err)
			logger.AppLogger().Warnf("%+v", err1)
			return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
				Message: err1.Error()}, err1
		}
		logger.AppLogger().Debugf("verify recved bSignedBtid succ, bSignedBtid:%v", bSignedBtid)

		//generate sign
		if device_ability.GetAbilityModel().SecurityChipSupport {
			b64SignedBtid, err = device.SignFromSecurityChip([]byte(btid))
			if err != nil {
				err1 := fmt.Errorf("%+v", err)
				logger.AppLogger().Warnf("%+v", err1)
				return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
					Message: err1.Error()}, err1
			}

		} else {
			pri, err := encwrapper.GetPrivateKey(string(device.GetDevicePriKey()))
			if err != nil {
				err1 := fmt.Errorf("%+v", err)
				logger.AppLogger().Warnf("%+v", err1)
				return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
					Message: err1.Error()}, err1
			}
			d, err := encwrapper.Sign(pri, []byte(btid))
			if err != nil {
				err1 := fmt.Errorf("%+v", err)
				logger.AppLogger().Warnf("%+v", err1)
				return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
					Message: err1.Error()}, err1
			}
			b64SignedBtid = encoding.Base64Encode(d)
		}
	}

	results := &dtopair.PubKeyExchangeRsp{BoxPubKey: string(device.GetDevicePubKey()),
		SignedBtid: b64SignedBtid}

	rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		Message: "OK",
		Results: results}
	logger.AppLogger().Debugf("ServicePubKeyExchange, rsp:%+v, results:%+v", rsp, rsp.Results)
	return rsp, nil
}
