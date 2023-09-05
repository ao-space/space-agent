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

package create

import (
	"agent/biz/model/dto/did/document"
	"agent/biz/model/dto/pair"
)

type CreateReq struct {
	ClientUuid           string                         `json:"clientUuid" binding:"required"`         // 客户端 uuid
	ClientPhoneModel     string                         `json:"clientPhoneModel" binding:"required"`   // 客户端手机型号.
	SpaceName            string                         `json:"spaceName" binding:"required"`          // 空间名称
	Password             string                         `json:"password" binding:"required"`           // 空间密码
	EnableInternetAccess bool                           `json:"enableInternetAccess"`                  // 是否启用互联网通道
	PlatformApiBase      string                         `json:"platformApiBase"`                       // Platform api url setting, e.g. "https://ao.space"`
	VerifyMethod         []*document.VerificationMethod `json:"verificationMethod" binding:"required"` // did 的验证方法
}

type CreateRsp struct {
	AgentToken           string          `json:"agentToken"`              // agent 接口的访问 token. 暂未启用.
	EnableInternetAccess bool            `json:"enableInternetAccess"`    //是否启用互联网通道
	ConnectedNetwork     []*pair.Network `json:"connectedNetwork"`        // 设备连接的网络情况
	SpaceUserInfo        interface{}     `json:"spaceUserInfo,omitempty"` // 空间信息
	DIDDoc               string          `json:"didDoc,omitempty"`        // did 文档
	EncryptedPriKeyBytes string          `json:"encryptedPriKeyBytes"`    // 加密的空间密码凭证对应的私钥密文
}
