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

package encwrapper

import (
	"agent/utils/logger"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func GetPublicKey(publicKey string) (*rsa.PublicKey, error) {
	// publicKey = strings.ReplaceAll(publicKey, "-----BEGIN PUBLIC KEY-----\n", "")

	logger.AppLogger().Debugf("getPublicKey, publicKey=%v", publicKey)

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, fmt.Errorf("failed pem.Decode, block:%+v", block)
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed x509.ParsePKIXPublicKey, err:%+v", err)
	}
	pub := publicKeyInterface.(*rsa.PublicKey)
	return pub, nil
}

func GetPrivateKey(privateKey string) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, fmt.Errorf("failed pem.Decode, block:%+v", block)
	}
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	// privateKeyInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed x509.ParseECPrivateKey, err:%+v", err)
	}
	pri := privateKeyInterface.(*rsa.PrivateKey)
	// pri := privateKeyInterface
	return pri, nil
}

func Verify(pub *rsa.PublicKey, sign, data []byte) error {
	hash := sha256.New()
	hash.Write(data)
	dataSum := hash.Sum(nil)
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, dataSum, sign)
	return err
}

func Sign(pri *rsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.New()
	hash.Write(data)
	dataSum := hash.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, pri, crypto.SHA256, dataSum)
	return signature, err
}
