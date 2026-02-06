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

package device

import (
	"agent/biz/model/device_ability"
	"agent/config"
	"agent/utils/unixsock/http"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
	gocomRsa "github.com/dungeonsnd/gocom/encrypt/rsa"
)

var devicePriKeyBytes []byte
var devicePublicKey []byte

func testGetPublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	logger.AppLogger().Debugf("testGetPublicKey, publicKey=%v", string(publicKey))
	block, _ := pem.Decode(publicKey)
	if block == nil {
		logger.AppLogger().Errorf("testGetPublicKey, failed pem.Decode, block:%+v", block) // 致命错误!
		return nil, fmt.Errorf("failed pem.Decode, block:%+v", block)
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.AppLogger().Errorf("testGetPublicKey, failed x509.ParsePKIXPublicKey, err:%+v", err) // 致命错误!
		return nil, fmt.Errorf("failed x509.ParsePKIXPublicKey, err:%+v", err)
	}
	pub := publicKeyInterface.(*rsa.PublicKey)
	return pub, nil
}

func testGetPrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	logger.AppLogger().Debugf("testGetPrivateKey")
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		logger.AppLogger().Errorf("testGetPrivateKey, failed pem.Decode, block:%+v", block) // 致命错误!
		return nil, fmt.Errorf("failed pem.Decode, block:%+v", block)
	}
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	// privateKeyInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.AppLogger().Errorf("testGetPrivateKey, failed x509.ParseECPrivateKey, err:%+v", err) // 致命错误!
		return nil, fmt.Errorf("failed x509.ParseECPrivateKey, err:%+v", err)
	}
	pri := privateKeyInterface.(*rsa.PrivateKey)
	// pri := privateKeyInterface
	return pri, nil
}

func VerifyKeys(devicePriKeyBytes []byte, devicePublicKey []byte) error {
	if len(devicePriKeyBytes) < 1 || len(devicePublicKey) < 1 {
		return fmt.Errorf("VerifyKeys, devicePriKeyBytes length:%v, devicePublicKey length:%v", len(devicePriKeyBytes), len(devicePublicKey))
	}

	_, err1 := testGetPublicKey(devicePublicKey)
	if err1 != nil {
		return err1
	}
	_, err2 := testGetPrivateKey(devicePriKeyBytes)
	if err2 != nil {
		return err2
	}
	return nil
}

func InitDeviceKey() {
	logger.AppLogger().Debugf("InitdeviceKey")
	if device_ability.GetAbilityModel().SecurityChipSupport {
		logger.AppLogger().Debugf("HasSecurityChip")
	} else {
		InitDeviceKeyNormal()
	}
}

func InitDeviceKeyNormal() {
	var err error
	err, devicePriKeyBytes, devicePublicKey = gocomRsa.ReadRsaKeys(config.Config.Box.BoxKey.RsaKeyFile, config.Config.Box.BoxKey.RsaPubKeyFile)
	if err != nil || VerifyKeys(devicePriKeyBytes, devicePublicKey) != nil {
		logger.AppLogger().Debugf("InitdeviceKey, ReadRsaKeys err: %v, devicePriKeyBytes length:%v, devicePublicKey length:%v",
			err, len(devicePriKeyBytes), len(devicePublicKey))

		for i := 0; i < 3; i++ {
			err := gocomRsa.GenRsaKey(config.Config.Box.BoxKey.RsaKeyFile, config.Config.Box.BoxKey.RsaPubKeyFile, 2048)
			if err != nil {
				logger.AppLogger().Errorf("InitdeviceKey, genRsaKey err=%v", err) // 致命错误!
			} else {
				logger.AppLogger().Debugf("InitdeviceKey, genRsaKey sucess")
				break
			}

			time.Sleep(time.Duration(1) * time.Second)
		}
	}
	err, devicePriKeyBytes, devicePublicKey = gocomRsa.ReadRsaKeys(config.Config.Box.BoxKey.RsaKeyFile, config.Config.Box.BoxKey.RsaPubKeyFile)
	if err != nil {
		logger.AppLogger().Errorf("InitdeviceKey, readRsaKey err=%v", err) // 致命错误!
	}

	logger.AppLogger().Debugf("InitdeviceKey, devicePriKeyBytes length=%v", len(devicePublicKey))
	logger.AppLogger().Debugf("InitdeviceKey, devicePublicKey length=%v", len(devicePublicKey))
	// logger.AppLogger().Debugf("InitdeviceKey, devicePublicKey=%v", devicePublicKey)
}

func GetDevicePriKey() []byte {
	logger.AppLogger().Debugf("GetDevicePriKey")
	if device_ability.GetAbilityModel().SecurityChipSupport {
		return nil
	} else {
		return devicePriKeyBytes
	}
}

func GetDevicePubKey() []byte {
	logger.AppLogger().Debugf("GetBoxPubKey")

	if device_ability.GetAbilityModel().SecurityChipSupport {
		logger.AppLogger().Debugf("GetBoxPubKey, 2gen board")
		pubkey, err := getFromSecurityChip(config.Config.Box.SecurityChipAgentSockAddr,
			"/security/v1/api/crypto/exportpubkey")
		if err != nil {
			logger.AppLogger().Errorf("failed getBoxPubKeyFromSecurityChip, err:%v", err)
			return nil
		}
		return []byte(pubkey)
	} else {
		logger.AppLogger().Debugf("GetBoxPubKey, 1gen board")
		return devicePublicKey
	}
}

func SignFromSecurityChip(data []byte) (string, error) {
	return postFromSecurityChip(config.Config.Box.SecurityChipAgentSockAddr,
		"/security/v1/api/crypto/sign",
		data)
}

func postFromSecurityChip(sockAddr, path string, data []byte) (string, error) {

	type Request struct {
		Input string `json:"input"`
	}

	type Result struct {
		Output string `json:"output"`
	}

	type Response struct {
		RequestId string `json:"requestId, omitempty"`
		Message   string `json:"message, omitempty"`
		Results   Result `json:"results, omitempty"`
	}

	parms := &Request{Input: encoding.Base64Encode(data)}
	var rsp Response
	_, _, _, err := http.PostJsonWithHeadersByUnixSock(sockAddr, path, parms, nil, &rsp)
	if err != nil {
		// fmt.Printf("enc failed, err:%v \n", err)
		return "", err
	}
	return rsp.Results.Output, nil
}

func getFromSecurityChip(sockAddr string, path string) (string, error) {

	type Request struct {
		Input string `json:"input"`
	}

	type Result struct {
		Output string `json:"output"`
	}

	type Response struct {
		RequestId string `json:"requestId, omitempty"`
		Message   string `json:"message, omitempty"`
		Results   Result `json:"results, omitempty"`
	}
	var rspPubkey Response
	_, err := http.GetJsonWithHeadersByUnixSock(sockAddr, path, nil, nil, &rspPubkey)
	if err != nil {
		// fmt.Printf("exportpubkey failed, err:%v \n", err)
		return "", err
	}
	// fmt.Printf("exportpubkey result: %+v \n", rspPubkey.Results)
	return rspPubkey.Results.Output, nil
}

func GetDevicePubKeyFingerprint() (string, error) {
	return calRsaKeyFingerprint(string(GetDevicePubKey()))
}

/**
 * 计算公钥的指纹
 * 去掉公钥的首行和尾行的标识. 中间 base64 部分先base64解码, 然后md5, 最后hexstring.
 *
 * @author wenchao
 * @date 2021-10-12 22:55:10
 * @return String 指纹
 **/
func calRsaKeyFingerprint(pemKey string) (string, error) {
	res := ""

	// 这里应该可以用正则，但是不同语言可能需要依赖其他库，所以方便集成就自己做处理了。
	lineEnder := "\r\n"
	arr := strings.Split(pemKey, lineEnder)
	if len(arr) <= 1 {
		lineEnder = "\n"
		arr = strings.Split(pemKey, lineEnder)
	}
	if len(arr) <= 1 {
		lineEnder = "\r"
		arr = strings.Split(pemKey, lineEnder)
	}
	if len(arr) <= 3 {
		return res, errors.New("pem format error")
	}

	if len(arr) > 3 {
		for i := 1; i < len(arr)-1; i++ {
			if i != 1 {
				res += lineEnder
			}
			if len(arr[i]) > 0 && !strings.Contains(arr[i], "KEY") {
				res += arr[i]
			}
		}
	}

	d, err := encoding.Base64Decode(res)
	if err != nil {
		return res, fmt.Errorf("pem base64 decode error, err:%v", err)
	}

	t := md5.Sum(d)

	// fmt.Printf("t(%v):\n%v\n", len(t), t)

	return hex.EncodeToString(t[0:]), nil
}
