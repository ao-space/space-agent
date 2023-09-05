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

package jwt

import (
	"crypto"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"github.com/golang-jwt/jwt/v5"
)

type SigningMethodCustomSt struct {
	Name string
	Hash crypto.Hash
}

var (
	SigningMethodCustom *SigningMethodCustomSt
)

func init() {
	SigningMethodCustom = &SigningMethodCustomSt{"MethodCustom", crypto.SHA256}
	jwt.RegisterSigningMethod(SigningMethodCustom.Alg(), func() jwt.SigningMethod {
		return SigningMethodCustom
	})

}

func (m *SigningMethodCustomSt) Alg() string {
	return m.Name
}

// Verify implements token verification for the SigningMethod
// For this signing method, must be an *rsa.PublicKey structure.
func (m *SigningMethodCustomSt) Verify(signingString string, sig []byte, key interface{}) error {
	return verifyFromSecurityChip([]byte(signingString), string(sig))
}

// Sign implements token signing for the SigningMethod
// For this signing method, must be an *rsa.PrivateKey structure.
func (m *SigningMethodCustomSt) Sign(signingString string, key interface{}) ([]byte, error) {
	ret, err := signFromSecurityChip([]byte(signingString))
	return []byte(ret), err
}

func signFromSecurityChip(data []byte) (string, error) {
	url := getSignUrl()

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
	err := PostHttpRequest(url, parms, nil, &rsp)
	if err != nil {
		// fmt.Printf("enc failed, err:%v \n", err)
		return "", err
	}
	return rsp.Results.Output, nil
}

func verifyFromSecurityChip(data []byte, signature string) error {
	url := getVerifyUrl()

	type Request struct {
		Input     string `json:"input"`
		Signature string `json:"signature"`
	}

	type Result struct {
		Output string `json:"output"`
	}

	type Response struct {
		RequestId string `json:"requestId, omitempty"`
		Message   string `json:"message, omitempty"`
		Results   Result `json:"results, omitempty"`
	}

	parms := &Request{Input: encoding.Base64Encode(data), Signature: encoding.Base64Encode([]byte(signature))}
	var rsp Response
	err := PostHttpRequest(url, parms, nil, &rsp)
	if err != nil {
		// fmt.Printf("enc failed, err:%v \n", err)
		return err
	}
	return nil
}
