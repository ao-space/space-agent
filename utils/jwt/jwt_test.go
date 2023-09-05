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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"github.com/dungeonsnd/gocom/encrypt/random"
	gocomRsa "github.com/dungeonsnd/gocom/encrypt/rsa"
)

func getPublicKey(publicKey []byte) (*rsa.PublicKey, error) {

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

func getPrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, fmt.Errorf("failed pem.Decode, block:%+v", block)
	}
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed x509.ParseECPrivateKey, err:%+v", err)
	}
	pri := privateKeyInterface.(*rsa.PrivateKey)
	return pri, nil
}

func TestJWT(t *testing.T) {
	expiredAt := time.Now().Add(24 * time.Hour * 365 * 3)

	pri, pub, err := gocomRsa.GenRsaKeyToString(2048)
	if err != nil {
		t.Errorf("failed GenRsaKey, err:%v", err)
	}

	priKey, err := getPrivateKey([]byte(pri))
	if err != nil {
		t.Errorf("failed getPrivateKey, err:%v", err)
	}

	pubKey, err := getPublicKey([]byte(pub))
	if err != nil {
		t.Errorf("failed getGetPublicKey, err:%v", err)
	}

	boxUuid := random.GenUUID()
	clientUuid := random.GenUUID()
	tokenType := "BIND_API_TOKEN"

	jwtToken, err := GenerateJWT(boxUuid, "access_token", []string{clientUuid}, expiredAt, nil, map[string]string{"tokenType": tokenType}, priKey)
	if err != nil {
		t.Errorf("failed GenerateJWT, err:%v", err)
	}
	t.Logf("jwtToken:%+v", jwtToken)

	issuer, subject, audience, m1, err := ParseJwt(jwtToken, pubKey)
	if err != nil {
		t.Errorf("failed ParseJwt, err:%v", err)
	}
	t.Logf("issuer:%+v, subject:%+v, audience:%+v, m1:%+v", issuer, subject, audience, m1)

	if issuer != boxUuid {
		t.Errorf("issuer{%v} != boxUuid{%v}", string(issuer), boxUuid)
	}
	if len(audience) > 0 && audience[0] != clientUuid {
		t.Errorf("audience[0]{%v} != clientUuid{%v}", audience[0], clientUuid)
	}
}

func TestJWTCustom(t *testing.T) {
	expiredAt := time.Now().Add(24 * time.Hour * 365 * 3)

	boxUuid := random.GenUUID()
	clientUuid := random.GenUUID()
	tokenType := "BIND_API_TOKEN"

	jwtToken, err := GenerateJWT(boxUuid, "access_token", []string{clientUuid}, expiredAt, nil, map[string]string{"tokenType": tokenType}, nil)
	if err != nil {
		t.Errorf("failed GenerateJWT, err:%v", err)
	}
	t.Logf("jwtToken:%+v", jwtToken)

	issuer, subject, audience, m1, err := ParseJwt(jwtToken, nil)
	if err != nil {
		t.Errorf("failed ParseJwt, err:%v", err)
	}
	t.Logf("issuer:%+v, subject:%+v, audience:%+v, m1:%+v", issuer, subject, audience, m1)

	if issuer != boxUuid {
		t.Errorf("issuer{%v} != boxUuid{%v}", string(issuer), boxUuid)
	}
	if len(audience) > 0 && audience[0] != clientUuid {
		t.Errorf("audience[0]{%v} != clientUuid{%v}", audience[0], clientUuid)
	}
}

func TestJWTEncryptFileds(t *testing.T) {
	expiredAt := time.Now().Add(24 * time.Hour * 365 * 5)

	pri, pub, err := gocomRsa.GenRsaKeyToString(2048)
	if err != nil {
		t.Errorf("failed GenRsaKey, err:%v", err)
	}

	priKey, err := getPrivateKey([]byte(pri))
	if err != nil {
		t.Errorf("failed getPrivateKey, err:%v", err)
	}

	pubKey, err := getPublicKey([]byte(pub))
	if err != nil {
		t.Errorf("failed getGetPublicKey, err:%v", err)
	}

	boxUuid := random.GenUUID()
	clientUuid := random.GenUUID()
	encBoxUuid, err := gocomRsa.RsaEncrypt([]byte(boxUuid), []byte(pub))
	if err != nil {
		t.Errorf("failed RsaEncrypt, err:%v", err)
	}

	encClientUuid, err := gocomRsa.RsaEncrypt([]byte(clientUuid), []byte(pub))
	if err != nil {
		t.Errorf("failed RsaEncrypt, err:%v", err)
	}

	m := map[string]string{}
	m["boxUuid"] = encoding.Base64Encode(encBoxUuid)
	m["clientUuid"] = encoding.Base64Encode(encClientUuid)
	jwtToken, err := GenerateJWT("box", "access_token", []string{"client"}, expiredAt, m, nil, priKey)
	if err != nil {
		t.Errorf("failed GenerateJWT, err:%v", err)
	}
	t.Logf("jwtToken:%+v", jwtToken)

	issuer, subject, audience, m1, err := ParseJwt(jwtToken, pubKey)
	if err != nil {
		t.Errorf("failed ParseJwt, err:%v", err)
	}
	t.Logf("issuer:%+v, subject:%+v, audience:%+v", issuer, subject, audience)

	decodeBoxUuid, err := encoding.Base64Decode(m1["boxUuid"])
	if err != nil {
		t.Errorf("failed Base64Decode, err:%v", err)
	}
	decodeClientUuid, err := encoding.Base64Decode(m1["clientUuid"])
	if err != nil {
		t.Errorf("failed Base64Decode, err:%v", err)
	}
	decBoxUuid, err := gocomRsa.RsaDecrypt([]byte(decodeBoxUuid), []byte(pri))
	if err != nil {
		t.Errorf("failed RsaDecrypt, err:%v", err)
	}
	decClientUuid, err := gocomRsa.RsaDecrypt([]byte(decodeClientUuid), []byte(pri))
	if err != nil {
		t.Errorf("failed RsaDecrypt, err:%v", err)
	}
	t.Logf("decBoxUuid:%+v, decClientUuid:%+v", string(decBoxUuid), string(decClientUuid))

	if string(decBoxUuid) != boxUuid {
		t.Errorf("decBoxUuid{%v} != boxUuid{%v}", string(decBoxUuid), boxUuid)
	}
	if string(decClientUuid) != clientUuid {
		t.Errorf("decClientUuid{%v} != clientUuid{%v}", string(decClientUuid), clientUuid)
	}
}
