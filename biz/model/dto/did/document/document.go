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

package document

const (
	CredentialTypePassword = "password"
)

type VerificationMethod struct {
	ID           string `json:"id" form:"id" binding:"required"`
	Type         string `json:"type" form:"type" binding:"required"`
	PublicKeyPem string `json:"publicKeyPem" form:"publicKeyPem" binding:"required"`
}

type GetDocumentReq struct {
	DID  string `json:"did" form:"did"`
	AOID string `json:"aoId" form:"aoId"`
}

type GetDocumentRsp struct {
	DIDDoc               string `json:"didDoc,omitempty"`
	EncryptedPriKeyBytes string `json:"encryptedPriKeyBytes"`
}
