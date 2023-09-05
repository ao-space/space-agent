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

package did

import (
	"agent/biz/model/did/leveldb"
	"agent/biz/model/dto/did/document"
	"testing"
)

func TestCreateDocument(t *testing.T) {
	if err := leveldb.OpenDB(); err != nil {
		panic(err)
	}

	t.Logf("\n$$$$ CreateDocument\n")
	aoId := "aoId-1"
	oldPassword := "123456"
	newPassword := "111111"
	ID := ":AAAHtMWCPnvz2q5ONvw="
	keyType := "RsaVerificationKey2018"
	publicKeyPemClient := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnN5jap7CGcqYURbLDVUa\nLc9kMxOyCMEykfwbQKXvTkPMkR9tKZmq8EqfG2d2OyUpF1TIfqHK7Q6d33yD02oO\nBTXZw1Ijkfxvu0KwG2zLV02FTuwZzgYa/AaP5iRZDx5GwTk/YFw+NTqT8Gf29a/L\n/ItcCfsEFLr3zMDXUcU9A7rBEy5ncva6RLNpXawegFGlCZa5+Gah8voKl8ZGpIgt\nlSc1IdnbPbBCYYlUATWLCLeYl+Q9/LslbpkFtdR+4M8vU7G1H+AQZ5fr2E9qX36I\nzcnchDmKq5bkbWQ9GJeZKqZTkhtCPBy4cphM8fHtZuoh1fA3VfF01N4KHT2bUdtp\nJwIDAQAB\n-----END PUBLIC KEY-----"
	verificationMethod := &document.VerificationMethod{ID: ID, Type: keyType, PublicKeyPem: publicKeyPemClient}
	verificationMethods := []*document.VerificationMethod{verificationMethod}
	_, didDocBytes, did, err := CreateDocument(aoId, oldPassword, verificationMethods)
	if err != nil {
		panic(err)
	}

	t.Logf("\ndidDocBytes:%+v\n", string(didDocBytes))
	t.Logf("\ndid:%+v\n", did)

	t.Logf("\n$$$$ UpdateDocumentOfPasswordVerficationByDid\n")
	err = UpdatePasswordKey(did, aoId, oldPassword, newPassword)
	if err != nil {
		panic(err)
	}

	t.Logf("\n$$$$ GetDocumentFromFile\n")
	didDocBytes, err = GetDocumentFromFile(did)
	if err != nil {
		panic(err)
	}
	t.Logf("\ndidDocBytes:%+v\n", string(didDocBytes))
	t.Logf("\ndid:%+v\n", did)
}
