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
 * @Date: 2021-10-16 10:23:56
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-17 17:17:33
 * @Description:
 */
package pair

type PubKeyExchangeReq struct {
	ClientPubKey string `json:"clientPubKey"` // 客户端公钥.
	ClientPriKey string `json:"clientPriKey"` // 客户端私钥.
	SignedBtid   string `json:"signedBtid"`   // 用客户端私钥签名btid. 盒子端用客户端公钥进行验证.
}

type PubKeyExchangeRsp struct {
	BoxPubKey  string `json:"boxPubKey"`  // 盒子端公钥.
	SignedBtid string `json:"signedBtid"` // 用盒子端端私钥签名btid. 客户端用盒子端公钥进行验证.
}

type KeyExchangeReq struct {
	ClientPreSecret string `json:"clientPreSecret"` // 必填,客户端对称密钥种子
	EncBtid         string `json:"encBtid"`         // 必填,使用盒子端公钥加密btid后进行base64得到的字符串
}
type KeyExchangeRsp struct {
	SharedSecret string `json:"sharedSecret"` // 对称密钥.
	Iv           string `json:"iv"`           // iv.
}

type PairingReq struct {
	ClientUuid       string `json:"clientUuid"`       // 客户端唯一id.
	ClientPubKey     string `json:"clientPubKey"`     // 客户端公钥.
	ClientPriKey     string `json:"clientPriKey"`     // 客户端私钥.
	ClientPhoneModel string `json:"clientPhoneModel"` // 客户端手机型号.
}

// type PairingBoxInfo struct {
// 	BoxUuid    string `json:"boxUuid"`    // 盒子端唯一id.
// 	BoxPubKey  string `json:"boxPubKey"`  // 盒子端公钥.
// 	AuthKey    string `json:"authKey"`    // 盒子端授权给客户端的凭证.
// 	RegKey     string `json:"regKey"`     // 平台端授权给客户端的凭证.
// 	UserDomain string `json:"userDomain"` // 盒子的唯一性域名(当前是盒子自动生成).
// 	BoxName    string `json:"boxName"`    // 盒子端名称
// }

// type PairingBoxInfoEnc struct {
// 	Data string `json:"data"` // 加密数据
// }

// c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "result": paringBoxInfo})
// type PairingBoxInfoResult struct {
// 	Results string `json:"results,omitempty"`
// }

// type PairingBoxInfoRspTest struct {
// 	Code    int             `json:"code"`
// 	Message string          `json:"message"`
// 	Result  *PairingBoxInfo `json:"result,omitempty"`
// }
