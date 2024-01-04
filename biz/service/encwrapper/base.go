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
 * @Date: 2021-12-13 16:43:09
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-24 16:28:09
 * @Description:
 */
package encwrapper

import (
	"agent/biz/model/clientinfo"
	"agent/biz/model/dto"
	"agent/utils/crypto"
	"fmt"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
)

func Check() error {
	_, _, err := clientinfo.GetSharedSecret()
	if err != nil {
		return err
	}
	return nil
}

func Enc(results interface{}) (dto.BaseRspStr, error) {
	logger.AppLogger().Debugf("enc, plain results:%+v", results)
	if results == nil {
		rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr,
			Message: "OK"}
		return rsp, nil
	}

	key, iv, _ := clientinfo.GetSharedSecret()

	d, err := encoding.JsonEncode(results)
	if err != nil {
		err1 := fmt.Errorf("enc err:%v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err1.Error()}, nil
	}
	encData, err := crypto.EncryptByAesAndBase64(d, []byte(key), []byte(iv))
	if err != nil {
		err1 := fmt.Errorf("enc err:%v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err1.Error()}, nil
	}

	rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		Message: "OK",
		Results: encData}
	logger.AppLogger().Debugf("enc, results:%+v, rsp:%+v", string(d), rsp)
	logger.AccessLogger().Debugf("[Enc], results:%+v, rsp:%+v", string(d), rsp)
	return rsp, nil
}

func EncBytes(requestId string, results []byte) (dto.BaseRspStr, error) {
	logger.AppLogger().Debugf("EncBytes, plain len(results):%+v", len(results))
	key, iv, _ := clientinfo.GetSharedSecret()

	encData, err := crypto.EncryptByAesAndBase64(results, []byte(key), []byte(iv))
	if err != nil {
		err1 := fmt.Errorf("enc err:%v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err1.Error()}, nil
	}

	rsp := dto.BaseRspStr{Code: dto.AgentCodeOkStr,
		RequestId: requestId,
		Message:   "OK",
		Results:   encData}
	logger.AppLogger().Debugf("enc, results:%+v, rsp:%+v", results, rsp)
	return rsp, nil
}

func DecParam(cipher string) ([]byte, error) {
	if len(cipher) < 1 {
		return nil, fmt.Errorf("decParam, cipher len is 0")
	}
	key, iv, _ := clientinfo.GetSharedSecret()

	decData, err := crypto.DecryptByAesAndBase64(cipher, []byte(key), []byte(iv))
	if err != nil {
		return nil, err
	}
	return decData, nil
}

func Dec(params ...string) ([]string, error) {
	if len(params) < 1 {
		err1 := fmt.Errorf("dec, params=%v", params)
		logger.AppLogger().Warnf("dec, params=%v", params)
		return nil, err1
	}
	rt := make([]string, 0)
	for _, v := range params {
		d, err := DecParam(v)
		if err != nil {
			return nil, err
		}
		// logger.AppLogger().Debugf("dec, parm[%v]=%v", i, v)
		rt = append(rt, string(d))
	}
	return rt, nil
}
