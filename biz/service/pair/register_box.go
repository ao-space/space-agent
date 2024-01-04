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
 * @Date: 2021-12-10 13:33:26
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-21 16:39:40
 * @Description:
 */

package pair

import (
	"agent/biz/model/device"
	"agent/biz/model/device_ability"
	"agent/biz/service/encwrapper"
	"agent/config"
	"agent/utils"
	"fmt"
	"net/http"
	"time"

	"agent/utils/logger"

	utilshttp "agent/utils/network/http"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"github.com/dungeonsnd/gocom/encrypt/random"
)

// 向平台注册盒子
func ServiceRegisterBox() error {
	logger.AppLogger().Debugf("ServiceRegisterBox")

	//先获取 box-reg-key
	if boxRegKeyInfo, err := GetDeviceRegKey(""); err != nil {
		logger.AppLogger().Warnf("ServiceRegisterBox, failed GetBoxRegKey, err:%+v", err)
		return err
	} else {
		logger.AppLogger().Debugf("ServiceRegisterBox, succ GetBoxRegKey, boxRegKeyInfo:%+v", boxRegKeyInfo)
		device.SetDeviceRegKey(boxRegKeyInfo.BoxRegKey, boxRegKeyInfo.ExpiresAt)
	}

	// 平台请求结构
	type registryStruct struct {
		BoxUUID string `json:"boxUUID"`
	}
	// 平台响应结构
	type registryRspStruct struct {
		BoxUUID       string `json:"boxUUID"`
		NetworkClient struct {
			ClientID  string `json:"clientId"`
			SecretKey string `json:"secretKey"`
		} `json:"networkClient"`
	}

	// 请求平台
	parms := &registryStruct{BoxUUID: device.GetDeviceInfo().BoxUuid}
	// url := config.Config.Platform.APIBase.Url + config.Config.Platform.RegistryBox.Path
	url := device.GetApiBaseUrl() + config.Config.Platform.RegistryBox.Path
	logger.AppLogger().Debugf("ServiceRegisterBox, v2, url:%+v, parms:%+v", url, parms)

	var headers = map[string]string{"Request-Id": random.GenUUID(), "Box-Reg-Key": device.GetDeviceInfo().BoxRegKey}
	var rsp registryRspStruct

	tryTotal := 3
	// var httpReq *http.Request
	var httpRsp *http.Response
	var body []byte
	var err1 error
	for i := 0; i < tryTotal; i++ {
		_, httpRsp, body, err1 = utilshttp.PostJsonWithHeaders(url, parms, headers, &rsp)
		if err1 != nil {
			// logger.AppLogger().Warnf("Failed PostJson, err:%v, @@httpReq:%+v, @@httpRsp:%+v, @@body:%v", err1, httpReq, httpRsp, string(body))
			if i == tryTotal-1 {
				return err1
			}
			time.Sleep(time.Second * 2)
			continue
		} else {
			break
		}
	}

	logger.AppLogger().Infof("ServiceRegisterBox, parms:%+v", parms)
	logger.AppLogger().Infof("ServiceRegisterBox, rsp:%+v", rsp)
	logger.AppLogger().Infof("ServiceRegisterBox, httpRsp:%+v", httpRsp)

	if httpRsp.StatusCode == http.StatusOK {
		// 保存盒子信息
		device.SetNetworkClient(&device.NetworkClientInfo{ClientID: rsp.NetworkClient.ClientID,
			SecretKey: rsp.NetworkClient.SecretKey})
	} else if httpRsp.StatusCode == http.StatusNotAcceptable {
		boxInfo := device.GetDeviceInfo()
		if len(boxInfo.BoxRegKey) < 1 {
			logger.AppLogger().Warnf("ServiceRegisterBox, boxInfo.BoxRegKey: %+v", boxInfo.BoxRegKey)
			logger.AppLogger().Warnf("ServiceRegisterBox, boxInfo.NetworkClient: %+v", boxInfo.NetworkClient)
			return fmt.Errorf("box uuid had already registered in platform. Plz reset first!")
		} else {
			logger.AppLogger().Infof("ServiceRegisterBox, using exist BoxInfo: %+v", boxInfo)
		}
	} else {
		return fmt.Errorf("httpRsp.StatusCode=%v, @@body:%v", httpRsp.StatusCode, string(body))
	}
	return nil
}

type BoxRegKeyInfo struct {
	ServiceID string `json:"serviceId"`
	BoxRegKey string `json:"boxRegKey"`
	ExpiresAt string `json:"expiresAt"`
}

func GetDeviceRegKey(apiBaseUrl string) (*BoxRegKeyInfo, error) {

	logger.AppLogger().Debugf("getBoxRegKey")

	// 平台请求结构
	type authStruct struct {
		BoxUUID    string   `json:"boxUUID"`
		ServiceIds []string `json:"serviceIds"`
		Sign       string   `json:"sign,omitempty"`
	}

	// 平台响应结构
	type authRspStruct struct {
		BoxUUID      string          `json:"boxUUID"`
		TokenResults []BoxRegKeyInfo `json:"tokenResults"`
	}

	logger.AppLogger().Debugf("getBoxRegKey, apiBaseUrl:%+v, boxInfo.ApiBaseUrl:%+v, config.Config.Platform.APIBase.Url:%+v",
		apiBaseUrl, device.GetDeviceInfo().ApiBaseUrl, config.Config.Platform.APIBase.Url)
	if len(apiBaseUrl) == 0 {
		apiBaseUrl = device.GetApiBaseUrl()
	}

	sign := ""
	// 生成签名
	signObj := &authStruct{BoxUUID: device.GetDeviceInfo().BoxUuid,
		ServiceIds: []string{"10001"}}
	signBytes, err := encoding.JsonEncode(signObj)
	if err != nil {
		err1 := fmt.Errorf("%+v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return nil, err1
	}
	logger.AppLogger().Debugf("GetBoxRegKey, signBytes:%v", string(signBytes))

	if device_ability.GetAbilityModel().SecurityChipSupport {
		logger.AppLogger().Debugf("GetBoxRegKey, SecurityChipSupport")
		sign, err = device.SignFromSecurityChip(signBytes)
		if err != nil {
			err1 := fmt.Errorf("%+v", err)
			logger.AppLogger().Warnf("%+v", err1)
			return nil, err1
		}
		logger.AppLogger().Debugf("GetBoxRegKey,  sign:%v", sign)

	} else {
		logger.AppLogger().Debugf("GetBoxRegKey, SecurityChipSupport==false")
		pri, err := encwrapper.GetPrivateKey(string(device.GetDevicePriKey()))
		if err != nil {
			err1 := fmt.Errorf("%+v", err)
			logger.AppLogger().Warnf("%+v", err1)
			return nil, err1
		}
		d, err := encwrapper.Sign(pri, signBytes)
		if err != nil {
			err1 := fmt.Errorf("%+v", err)
			logger.AppLogger().Warnf("%+v", err1)
			return nil, err1
		}
		sign = encoding.Base64Encode(d)
		logger.AppLogger().Debugf("GetBoxRegKey, SecurityChipSupport==false, sign:%v", sign)
	}

	// 请求平台
	parms := &authStruct{BoxUUID: device.GetDeviceInfo().BoxUuid,
		ServiceIds: []string{"10001"},
		Sign:       sign}
	url, _ := utils.JoinUrl(apiBaseUrl, config.Config.Platform.AuthBox.Path)

	logger.AppLogger().Debugf("getBoxRegKey, url:%+v, parms:%+v", url, parms)
	var headers = map[string]string{"Request-Id": random.GenUUID()}
	var rsp authRspStruct
	tryTotal := 3
	var httpReq *http.Request
	var httpRsp *http.Response
	var body []byte
	var err1 error
	for i := 0; i < tryTotal; i++ {
		httpReq, httpRsp, body, err1 = utilshttp.PostJsonWithHeaders(url, parms, headers, &rsp)
		if err1 != nil {
			logger.AppLogger().Warnf("Failed PostJson, err:%v, @@httpReq:%+v, @@httpRsp:%+v, @@body:%v", err1, httpReq, httpRsp, string(body))
			if i == tryTotal-1 {
				return nil, err1
			}
			time.Sleep(time.Second * 2)
			continue
		} else {
			break
		}
	}
	logger.AppLogger().Infof("getBoxRegKey, httpReq:%+v", httpReq)
	logger.AppLogger().Infof("getBoxRegKey, parms:%+v", parms)
	logger.AppLogger().Infof("getBoxRegKey, httpRsp:%+v", httpRsp)
	logger.AppLogger().Infof("getBoxRegKey, rsp:%+v", rsp)
	logger.AppLogger().Infof("getBoxRegKey, body:%v", string(body))

	if httpRsp.StatusCode == http.StatusOK {
		if len(rsp.TokenResults) == 0 || len(rsp.TokenResults) > 2 {
			return nil, fmt.Errorf("len(rsp.TokenResults)=%v", len(rsp.TokenResults))
		}
		// 保存盒子信息
		for _, token := range rsp.TokenResults {
			switch token.ServiceID {
			case "10001":
				return &token, nil
			case "10002":

			default:
				return nil, fmt.Errorf("invalid serviceId(%v)", token.ServiceID)
			}
		}
	} else {
		return nil, fmt.Errorf("httpRsp.StatusCode=%v, @@body:%v", httpRsp.StatusCode, string(body))
	}

	return nil, fmt.Errorf("failed to get box-reg-key")
}
