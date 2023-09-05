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
 * @Date: 2021-10-29 13:55:30
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-10 17:01:24
 * @Description:
 */
package crypto

import (
	"fmt"

	"github.com/dungeonsnd/gocom/encrypt/aes"
	"github.com/dungeonsnd/gocom/encrypt/encoding"
	gocomRsa "github.com/dungeonsnd/gocom/encrypt/rsa"
	"agent/utils/logger"
)

func DecryptByPriKey(key []byte, b64parm string) ([]byte, error) {
	b, err := encoding.Base64Decode(b64parm)
	if err != nil {
		logger.AppLogger().Warnf("Base64Decode %v failed. err：%v", b64parm, err)
		return nil, err
	}

	decParm, err := gocomRsa.RsaDecrypt(b, key)
	if err != nil {
		logger.AppLogger().Warnf("RsaDecrypt %v failed. err：%v", b64parm, err)
		return nil, err
	}

	return decParm, nil
}

func EncryptByPubKey(key []byte, data []byte) (string, error) {
	enc, err := gocomRsa.RsaEncrypt(data, key)
	if err != nil {
		logger.AppLogger().Warnf("RsaEncrypt failed. err：%v", err)
		return "", err
	}
	return encoding.Base64Encode(enc), nil
}

func EncryptByAesAndBase64(origData []byte, key []byte, iv []byte) (string, error) {
	encbuf, err := aes.AesEncryptByPkcs5Padding(origData, key, iv)
	if err != nil {
		err1 := fmt.Errorf("AesEncryptByPkcs5Padding failed, err:%+v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return "", err1
	}
	encoded := encoding.Base64Encode(encbuf)
	return encoded, nil
}

func DecryptByAesAndBase64(base64Data string, key []byte, iv []byte) ([]byte, error) {
	logger.AppLogger().Debugf("DecryptByAesAndBase64, base64Data=%v, key=%v, iv=%v", base64Data, string(key), string(iv))

	encData, err := encoding.Base64Decode(base64Data)
	if err != nil {
		err1 := fmt.Errorf("Base64Decode base64Data failed.  err:%+v, base64Data=%v", err, base64Data)
		logger.AppLogger().Warnf("%+v", err1)
		return nil, err1
	}

	decBuf, err := aes.AesDecryptByPkcs5Padding(encData, key, iv)
	if err != nil {
		err1 := fmt.Errorf("AesDecryptByPkcs5Padding failed. err:%+v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return nil, err1
	}
	return decBuf, nil
}
