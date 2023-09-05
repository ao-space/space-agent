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

package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func GenRsaKey(bits int) ([]byte, []byte, error) {
	// gen pri
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	derStream, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	priBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	priKeyBytes := pem.EncodeToMemory(priBlock)

	// gen pub
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubKeyBytes := pem.EncodeToMemory(publicBlock)

	return priKeyBytes, pubKeyBytes, nil
}

func GetPublicKey(publicKey []byte) (*rsa.PublicKey, error) {

	block, _ := pem.Decode(publicKey)
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

func GetPrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode(privateKey)
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

func GetRsaPubKeyByPriKeyBytes(priKeyBytes []byte) ([]byte, error) {
	privateKey, err := GetPrivateKey(priKeyBytes)
	if err != nil {
		return nil, err
	}

	// gen pub
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubKeyBytes := pem.EncodeToMemory(publicBlock)

	return pubKeyBytes, nil
}
