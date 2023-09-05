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
	aospacedid "agent/deps/did/aospace/did"
	"agent/utils/logger"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gibson042/canonicaljson-go"
)

const (
	queryKeyNameOfVersionTime    = "versionTime"
	queryKeyNameOfCredentialType = "credentialType"

	CredentialTypeDevice           = "device"
	CredentialTypeBinder           = "binder"
	CredentialTypePasswordOnDevice = "passwordondevice"
	CredentialTypePasswordOnBinder = "passwordonbinder"
)

func GetEncryptedPriKeyBytes(levelDBTrans *leveldb.Trans, aoId string) ([]byte, bool, error) {
	logger.AppLogger().Debugf("GetEncryptedPriKeyBytes, aoId:%v",
		aoId)
	return getEncryptedPrivatePasswordKey(levelDBTrans, aoId)
}

// 如果不需要开启事务，levelDBTrans 传入 nil.
func GetDocumentFromFile(levelDBTrans *leveldb.Trans, aoId, didStr string) ([]byte, error) {
	logger.AppLogger().Debugf("GetDocumentFromFile, didStr:%v, aoId:%v",
		didStr, aoId)

	if len(aoId) < 1 && len(didStr) < 1 {
		return nil, fmt.Errorf("did(%v) and aoId(%v) can't both be empty", didStr, aoId)
	}
	if len(aoId) < 1 {
		aoIdFound, found, err := GetAoIdByDid(levelDBTrans, didStr)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, fmt.Errorf("aoId not found by did(%v)", didStr)
		}

		aoId = aoIdFound
	}
	if len(didStr) < 1 {
		didStrFound, found, err := GetDidByAoId(levelDBTrans, aoId)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, fmt.Errorf("did not found by aoId(%v)", aoId)
		}

		didStr = didStrFound
	}

	return getDidDoc(levelDBTrans, didStr)
}

func UpdatePasswordKey(levelDBTrans *leveldb.Trans, did, aoId, oldPassword, newPassword string) error {
	logger.AppLogger().Debugf("UpdatePasswordKey, did:%v, aoId:%v",
		did, aoId)

	if len(aoId) < 1 {
		if len(did) < 1 {
			return fmt.Errorf("did(%v) and aoId(%v) can't be empty", did, aoId)
		}

		aoIdFound, found, err := GetAoIdByDid(levelDBTrans, did)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("aoId not found by did(%v)", did)
		}

		aoId = aoIdFound
	}

	err := updatePasswordKey(levelDBTrans, aoId, oldPassword, newPassword)
	if err != nil {
		return err
	}
	return nil
}

func ResetPasswordVerficationMethod(levelDBTrans *leveldb.Trans, didStr, aoId, newPassword string) ([]byte, string, error) {
	logger.AppLogger().Debugf("ResetPasswordVerficationMethod, didStr:%v, aoId:%v",
		didStr, aoId)

	timeNow := time.Now().UTC().Format(time.RFC3339)

	if len(aoId) < 1 && len(didStr) < 1 {
		return nil, "", fmt.Errorf("did(%v) and aoId(%v) can't both be empty", didStr, aoId)
	}
	if len(aoId) < 1 {
		aoIdFound, found, err := GetAoIdByDid(levelDBTrans, didStr)
		if err != nil {
			return nil, "", err
		}
		if !found {
			return nil, "", fmt.Errorf("aoId not found by did(%v)", didStr)
		}

		aoId = aoIdFound
	}
	if len(didStr) < 1 {
		didStrFound, found, err := GetDidByAoId(levelDBTrans, aoId)
		if err != nil {
			return nil, "", err
		}
		if !found {
			return nil, "", fmt.Errorf("did not found by aoId(%v)", aoId)
		}

		didStr = didStrFound
	}

	doc, err := getDidDoc(levelDBTrans, didStr)
	if err != nil {
		return nil, "", err
	}

	didDoc := &aospacedid.Document{}
	err = json.Unmarshal(doc, didDoc)
	if err != nil {
		return nil, "", err
	}

	did, err := aospacedid.FromDocument(didDoc)
	if err != nil {
		return nil, "", err
	}

	// add CredentialTypePasswordOnSpace method
	fragment := "key-2"
	_, keyId, query, err := createPasswordOnSpaceMethod(levelDBTrans, fragment, aoId, newPassword, timeNow, did)
	if err != nil {
		return nil, "", err
	}
	logger.AppLogger().Debugf("AddNewVerificationMethod, query:%v, keyId:%v",
		query, keyId)

	// add CredentialTypePasswordOnBinder method
	// fragment = "key-3"
	// keyId, query, found, err = createPasswordOnBinderMethod(fragment, timeNow, did, verificationMethods)
	// if err != nil {
	// 	return nil, "", err
	// }
	// logger.AppLogger().Debugf("AddNewVerificationMethod, query:%v, keyId:%v, found:%v",
	// 	query, keyId, found)

	// create new document
	didDoc = did.Document(true)
	didDocBytes, err := canonicaljson.Marshal(didDoc)
	if err != nil {
		return nil, "", fmt.Errorf("MarshalIndent err:%v", err)
	}

	if err := saveDidDoc(levelDBTrans, did.DID(), didDocBytes); err != nil {
		return nil, "", fmt.Errorf("saveDidDoc err:%v", err)
	}

	return didDocBytes, did.DID(), nil
}

func CreateDocument(levelDBTrans *leveldb.Trans, aoId, password string, verificationMethods []*document.VerificationMethod) ([]byte, []byte, string, error) {
	logger.AppLogger().Debugf("CreateDocument, aoId:%v, verificationMethods:%v",
		aoId, verificationMethods)

	timeNow := time.Now().UTC().Format(time.RFC3339)

	if len(verificationMethods) < 1 {
		return nil, nil, "", fmt.Errorf("verificationMethods is empty")
	}

	did, err := aospacedid.NewIdentifier()
	if err != nil {
		return nil, nil, "", fmt.Errorf("NewIdentifier err:%v", err)
	}

	firstVerificationMethod := make([]string, 0)
	secondVerificationMethod := make([]string, 0)

	// add device method
	fragment := "key-0"
	keyId, query, err := createDeviceMethod(levelDBTrans, fragment, aoId, timeNow, did)
	if err != nil {
		return nil, nil, "", fmt.Errorf("NewIdentifier err:%v", err)
	}
	firstVerificationMethod = append(firstVerificationMethod, fragment)
	logger.AppLogger().Debugf("AddNewVerificationMethod, query:%v, keyId:%v",
		query, keyId)

	// add binder method
	fragment = "key-1"
	keyId, query, err = createBinderMethod(fragment, timeNow, did, verificationMethods)
	if err != nil {
		return nil, nil, "", err
	}
	firstVerificationMethod = append(firstVerificationMethod, fragment)
	logger.AppLogger().Debugf("AddNewVerificationMethod, query:%v, keyId:%v",
		query, keyId)

	// add CredentialTypePasswordOnSpace method
	fragment = "key-2"
	encryptedPriKeyBytes, keyId, query, err := createPasswordOnSpaceMethod(levelDBTrans, fragment, aoId, password, timeNow, did)
	if err != nil {
		return nil, nil, "", err
	}
	secondVerificationMethod = append(secondVerificationMethod, fragment)
	logger.AppLogger().Debugf("AddNewVerificationMethod, query:%v, keyId:%v",
		query, keyId)

	// add CredentialTypePasswordOnBinder method
	fragment = "key-3"
	keyId, query, found, err := createPasswordOnBinderMethod(fragment, timeNow, did, verificationMethods)
	if err != nil {
		return nil, nil, "", err
	}
	if found {
		secondVerificationMethod = append(secondVerificationMethod, fragment)
		logger.AppLogger().Debugf("AddNewVerificationMethod, query:%v, keyId:%v",
			query, keyId)
	}

	// add mutisig
	did.AddNewVerificationMethodOfMultisig(firstVerificationMethod, secondVerificationMethod)
	logger.AppLogger().Debugf("AddNewVerificationMethodOfMultisig")

	// add CapabilityInvocation
	did.AddNewCapabilityInvocation()
	logger.AppLogger().Debugf("AddNewCapabilityInvocation")

	didDoc := did.Document(true)
	logger.AppLogger().Debugf("Document")
	// json 规范化
	// https://www.rfc-editor.org/rfc/rfc8785
	// https://github.com/gibson042/canonicaljson-go
	// https://github.com/cyberphone/json-canonicalization
	// https://github.com/oyamist/merkle-json
	// https://github.com/gowebpki/jcs
	// didDocBytes, err := json.MarshalIndent(didDoc, "", "  ")
	didDocBytes, err := canonicaljson.Marshal(didDoc)
	if err != nil {
		return nil, nil, "", fmt.Errorf("MarshalIndent err:%v", err)
	}
	logger.AppLogger().Debugf("Marshal")

	if err := saveDidDoc(levelDBTrans, did.DID(), didDocBytes); err != nil {
		return nil, nil, "", fmt.Errorf("saveDidDoc err:%v", err)
	}
	logger.AppLogger().Debugf("saveDidDoc")

	return encryptedPriKeyBytes, didDocBytes, did.DID(), nil
}
