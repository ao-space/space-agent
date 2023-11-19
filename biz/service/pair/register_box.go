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
	"agent/utils/logger"
	"fmt"
	"github.com/ao-space/platform-sdk-go/v2"
	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"time"
)

// 向平台注册盒子
func ServiceRegisterBox() error {
	logger.AppLogger().Debugf("ServiceRegisterBox")

	//先获取 box-reg-key
	var client *platform.Client
	var err error
	if client, err = GetSdkClientWithDeviceRegKey(""); err != nil {
		logger.AppLogger().Warnf("ServiceRegisterBox, failed GetBoxRegKey, err:%+v", err)
		return err
	} else {
		logger.AppLogger().Debugf("ServiceRegisterBox, succ GetSdkClientWithBoxRegKey, boxRegKeyInfo:%+v", client.TokenResults)
		device.SetDeviceRegKey(client.TokenResults.BoxRegKey, client.TokenResults.ExpiresAt.String())
	}

	tryTotal := 3
	var resp *platform.RegisterDeviceResponse
	for i := 0; i < tryTotal; i++ {
		resp, err = client.RegisterDevice()
		if err != nil {
			logger.AppLogger().Warnf("Failed PostJson, err:%v, @@resp:%+v", err, resp)
			if i == tryTotal-1 {
				return err
			}
			time.Sleep(time.Second * 2)
			continue
		} else {
			break
		}
	}

	logger.AppLogger().Infof("ServiceRegisterBox, rsp:%+v", resp)
	// 保存盒子信息
	device.SetNetworkClient(&device.NetworkClientInfo{ClientID: resp.NetWorkClient.ClientId,
		SecretKey: resp.NetWorkClient.SecretKey})
	return nil
}

func GetSdkClientWithDeviceRegKey(apiBaseUrl string) (*platform.Client, error) {

	logger.AppLogger().Debugf("getBoxRegKey")

	logger.AppLogger().Debugf("getBoxRegKey, apiBaseUrl:%+v, boxInfo.ApiBaseUrl:%+v, config.Config.Platform.APIBase.Url:%+v",
		apiBaseUrl, device.GetDeviceInfo().ApiBaseUrl, config.Config.Platform.APIBase.Url)
	if len(apiBaseUrl) == 0 {
		apiBaseUrl = device.GetApiBaseUrl()
	}

	client, err := platform.NewClientWithHost(apiBaseUrl, nil)
	if err != nil {
		logger.AppLogger().Errorf("%+v", err)
	}
	client.SetZapLogger(logger.AppLogger())

	type authStruct struct {
		BoxUUID    string   `json:"boxUUID"`
		ServiceIds []string `json:"serviceIds"`
		Sign       string   `json:"sign,omitempty"`
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

	tryTotal := 3
	var resp *platform.ObtainBoxRegKeyResponse
	var err1 error
	for i := 0; i < tryTotal; i++ {
		resp, err1 = client.ObtainBoxRegKey(&platform.ObtainBoxRegKeyRequest{
			BoxUUID:    device.GetDeviceInfo().BoxUuid,
			ServiceIds: []string{"10001"},
			Sign:       sign,
		})
		if err1 != nil {
			logger.AppLogger().Warnf("Failed PostJson, err:%v, @@resp:%+v", err1, resp)
			if i == tryTotal-1 {
				return nil, err1
			}
			time.Sleep(time.Second * 2)
			continue
		} else {
			break
		}
	}
	logger.AppLogger().Infof("getBoxRegKey, httpReq:%+v", resp)

	if len(resp.TokenResults) == 0 || len(resp.TokenResults) > 2 {
		return nil, fmt.Errorf("len(rsp.TokenResults)=%v", len(resp.TokenResults))
	}
	// 保存盒子信息
	for _, token := range resp.TokenResults {
		switch token.ServiceId {
		case "10001":
			return client, nil
		case "10002":

		default:
			return nil, fmt.Errorf("invalid serviceId(%v)", token.ServiceId)
		}
	}

	return nil, fmt.Errorf("failed to get box-reg-key")
}
