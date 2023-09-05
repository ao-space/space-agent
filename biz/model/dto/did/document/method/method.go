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

package method

import "agent/biz/model/dto/did/document"

type UpdateDocumentMethodReq struct {
	NewPassword  string                         `json:"newPassword" form:"newPassword"` // 重置空间密码时需要传入密码
	DID          string                         `json:"did" form:"did"`                 // did 和 aoId 需要至少传一个参数.
	AOID         string                         `json:"aoId" form:"aoId"`
	VerifyMethod []*document.VerificationMethod `json:"verificationMethod" form:"verificationMethod"` // 增加的验证方法
}

type UpdateDocumentMethodRsp struct {
	DIDDoc string `json:"didDoc,omitempty"`
	DID    string `json:"did,omitempty"`
}
