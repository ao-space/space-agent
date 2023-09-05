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
	"fmt"
	"strings"
)

func getVerificationMethodByCredentialType(verificationMethods []*document.VerificationMethod, credentialType string) (*document.VerificationMethod, bool) {
	for _, v := range verificationMethods {
		logger.AppLogger().Debugf("getVerificationMethodByCredentialType, v:%v,  find target: %v",
			v, queryKeyNameOfCredentialType+"="+credentialType)
		logger.AppLogger().Debugf("getVerificationMethodByCredentialType, strings.Index(v.ID, ?): %v",
			strings.Index(v.ID, "?"))
		if i := strings.Index(v.ID, "?"); i > 0 {
			query := v.ID[i+1:]
			logger.AppLogger().Debugf("getVerificationMethodByCredentialType, query: %v",
				query)
			if strings.Contains(query, queryKeyNameOfCredentialType+"="+credentialType) {
				return v, true
			}
		}
	}
	return nil, false
}

func createDeviceMethod(levelDBTrans *leveldb.Trans, fragment, aoId, timeNow string, did *aospacedid.Identifier) (string, string, error) {
	logger.AppLogger().Debugf("createDeviceMethod, aoId:%v, fragment:%v",
		aoId, fragment)

	_, spacePubKeyBytes, err := getSpaceKey(levelDBTrans, aoId)
	if err != nil {
		logger.AppLogger().Warnf("failed getSpaceKey, aoId:%v, err:%v",
			aoId, err)
		return "", "", fmt.Errorf("getSpaceKey err:%v", err)
	}

	keyType := aospacedid.KeyTypeRSA.String()
	query := fmt.Sprintf("%v=%v&%v=%v", queryKeyNameOfVersionTime, timeNow,
		queryKeyNameOfCredentialType, CredentialTypeDevice)
	keyId, err := did.AddNewVerificationMethod(keyType, string(spacePubKeyBytes), query, fragment)
	if err != nil {
		logger.AppLogger().Warnf("failed AddNewVerificationMethod, query:%v, err:%v",
			query, err)
		return "", "", fmt.Errorf("AddNewVerificationMethod err:%v", err)
	}

	return keyId, query, nil
}

func createBinderMethod(fragment, timeNow string, did *aospacedid.Identifier,
	verificationMethods []*document.VerificationMethod) (string, string, error) {
	logger.AppLogger().Debugf("createBinderMethod, fragment:%v, verificationMethods:%v",
		fragment, verificationMethods)

	verificationMethod, found := getVerificationMethodByCredentialType(verificationMethods, CredentialTypeBinder)
	if !found {
		return "", "", fmt.Errorf("getVerificationMethodByCredentialType, not found: %v=%v",
			queryKeyNameOfCredentialType, CredentialTypeBinder)
	}

	calKeyId := aospacedid.CalVerificationIdString(verificationMethod.PublicKeyPem)
	if !strings.Contains(verificationMethod.ID, calKeyId) {
		logger.AppLogger().Warnf("verificationMethod.ID:%v error, calculated KeyId:%v",
			verificationMethod.ID, calKeyId)
		// 校验
		return "", "", fmt.Errorf("verificationMethod.ID:%v error, calculated KeyId:%v",
			verificationMethod.ID, calKeyId)
	}

	keyType := verificationMethod.Type
	publicKeyPem := verificationMethod.PublicKeyPem
	query := fmt.Sprintf("%v=%v", queryKeyNameOfVersionTime, timeNow)
	if i := strings.Index(verificationMethod.ID, "?"); i > 0 {
		query += "&" + verificationMethod.ID[i+1:]
	}
	keyId, err := did.AddNewVerificationMethod(keyType, publicKeyPem, query, fragment)
	if err != nil {
		logger.AppLogger().Warnf("failed AddNewVerificationMethod, query:%v, err:%v",
			query, err)
		return "", "", fmt.Errorf("AddNewVerificationMethod err:%v", err)
	}
	return keyId, query, nil
}

func createPasswordOnSpaceMethod(levelDBTrans *leveldb.Trans, fragment, aoId, password, timeNow string,
	did *aospacedid.Identifier) ([]byte, string, string, error) {
	logger.AppLogger().Debugf("createPasswordOnSpaceMethod, fragment:%v, aoId:%v",
		fragment, aoId)

	queryType := fmt.Sprintf("%v=%v", queryKeyNameOfCredentialType, CredentialTypePasswordOnDevice)
	cnt, err := did.DeleteVerificationMethodOfQuery(queryType)
	if err != nil {
		return nil, "", "", err
	}
	logger.AppLogger().Debugf("createPasswordOnSpaceMethod, RemoveVerificationMethodByFragment cnt:%v",
		cnt)

	err = deletePasswordKey(levelDBTrans, aoId)
	if err != nil {
		return nil, "", "", err
	}
	encryptedPriKeyBytes, _, passwordPubKeyBytes, err := getPasswordKey(levelDBTrans, aoId, password)
	if err != nil {
		logger.AppLogger().Warnf("failed getPasswordKey, aoId:%v, err:%v",
			aoId, err)
		return nil, "", "", fmt.Errorf("getPasswordKey err:%v", err)
	}
	logger.AppLogger().Debugf("createPasswordOnSpaceMethod, getPasswordKey, passwordPubKeyBytes :%v",
		string(passwordPubKeyBytes))

	keyType := aospacedid.KeyTypeRSA.String()
	query := fmt.Sprintf("%v=%v&%v=%v", queryKeyNameOfVersionTime, timeNow,
		queryKeyNameOfCredentialType, CredentialTypePasswordOnDevice)
	keyId, err := did.AddNewVerificationMethod(keyType, string(passwordPubKeyBytes), query, fragment)
	if err != nil {
		logger.AppLogger().Warnf("failed AddNewVerificationMethod, query:%v, err:%v",
			query, err)
		return nil, "", "", fmt.Errorf("AddNewVerificationMethod err:%v", err)
	}
	return encryptedPriKeyBytes, keyId, query, nil
}

func createPasswordOnBinderMethod(fragment, timeNow string, did *aospacedid.Identifier, verificationMethods []*document.VerificationMethod) (string, string, bool, error) {
	logger.AppLogger().Debugf("createPasswordOnBinderMethod, fragment:%v, verificationMethods:%v",
		fragment, verificationMethods)

	verificationMethod, found := getVerificationMethodByCredentialType(verificationMethods, CredentialTypePasswordOnBinder)
	if !found {
		return "", "", found, nil
	} else {

		calKeyId := aospacedid.CalVerificationIdString(verificationMethod.PublicKeyPem)
		if !strings.Contains(verificationMethod.ID, calKeyId) {
			logger.AppLogger().Warnf("verificationMethod.ID:%v error, calculated KeyId:%v",
				verificationMethod.ID, calKeyId)
			// 校验
			return "", "", found, fmt.Errorf("verificationMethod.ID:%v error, calculated KeyId:%v",
				verificationMethod.ID, calKeyId)
		}

		queryType := fmt.Sprintf("%v=%v", queryKeyNameOfCredentialType, CredentialTypePasswordOnBinder)
		cnt, err := did.DeleteVerificationMethodOfQuery(queryType)
		if err != nil {
			return "", "", found, err
		}
		logger.AppLogger().Debugf("createPasswordOnBinderMethod, RemoveVerificationMethodByFragment cnt:%v",
			cnt)

		keyType := verificationMethod.Type
		publicKeyPem := verificationMethod.PublicKeyPem
		query := fmt.Sprintf("%v=%v", queryKeyNameOfVersionTime, timeNow)
		if i := strings.Index(verificationMethod.ID, "?"); i > 0 {
			query += "&" + verificationMethod.ID[i+1:]
		}
		keyId, err := did.AddNewVerificationMethod(keyType, publicKeyPem, query, fragment)
		if err != nil {
			logger.AppLogger().Warnf("failed AddNewVerificationMethod, query:%v, err:%v",
				query, err)
			return "", "", found, fmt.Errorf("AddNewVerificationMethod err:%v", err)
		}
		logger.AppLogger().Debugf("AddNewVerificationMethod, query:%v, keyId:%v",
			query, keyId)
		return keyId, query, found, nil
	}
}
