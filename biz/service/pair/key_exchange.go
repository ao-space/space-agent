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
 * @LastEditTime: 2022-01-29 09:17:31
 * @Description:
 */

package pair

import (
	"agent/biz/model/clientinfo"
	"agent/biz/model/dto"
	dtopair "agent/biz/model/dto/pair"
	"agent/utils/crypto"
	"fmt"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"github.com/dungeonsnd/gocom/encrypt/random"
)

func ServiceKeyExchange(req *dtopair.KeyExchangeReq) (dto.BaseRspStr, error) {

	if !clientinfo.ClientExchangePubKeyExchanged() {
		err1 := fmt.Errorf("no public key exchanged")
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err1.Error()}, err1
	}
	if len(req.ClientPreSecret) > 1024 {
		err1 := fmt.Errorf("ClientPreSecret length error")
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err1.Error()}, err1
	}

	key := random.RandMixHashHex()[:32]

	//clientinfo.SaveClientExchangeKey() // https://pm.eulix.xyz/bug-view-645.html
	err := clientinfo.SetSharedSecret(key)
	if err != nil {
		err1 := fmt.Errorf("failed SetSharedSecret, err:%v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeBadReqStr,
			Message: err.Error()}, err1
	}

	sharedSecret, rawIv, _ := clientinfo.GetSharedSecret()

	keyStr, err := crypto.EncryptByPubKey([]byte(clientinfo.GetClientExchangePubKey()), []byte(sharedSecret))
	if err != nil {
		err1 := fmt.Errorf("failed EncryptByPubKey, err:%v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err.Error()}, err1
	}

	sharedIv := encoding.Base64Encode([]byte(rawIv))
	ivStr, err := crypto.EncryptByPubKey([]byte(clientinfo.GetClientExchangePubKey()), []byte(sharedIv))
	if err != nil {
		err1 := fmt.Errorf("failed EncryptByPubKey, err:%v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err.Error()}, err1
	}

	results := &dtopair.KeyExchangeRsp{SharedSecret: keyStr, Iv: ivStr}
	rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr, Message: "OK", Results: results}
	logger.AppLogger().Debugf("ServiceKeyExchange, rsp:%+v, results:%+v", rsp, results)
	return rsp, nil
}
