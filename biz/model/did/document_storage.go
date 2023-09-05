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
	"agent/deps/did/aospace/rsa"
	"agent/utils/logger"
	cryptosha256 "crypto/sha256"
	"fmt"
	"strings"

	"github.com/dungeonsnd/gocom/encrypt/aes"
	"github.com/dungeonsnd/gocom/encrypt/hash/sha256"
	"github.com/dungeonsnd/gocom/encrypt/random"
	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSource            = `:4i4x{%3^>D/DT0Qa?BN!IXI1I,w&Vgjy@t)np2k)rrsKJ4^F6T^fb$N?L=kI8(FiEi-A|XnN[XyId-B%(:uoWPip9e`
	passwordKeyVersionLen = 4
	passwordKeyIVLen      = 16
)

var (
	preGeneratedSpaceKeyPri    []byte
	preGeneratedSpaceKeyPub    []byte
	preGeneratedPasswordKeyPri []byte
	preGeneratedPasswordKeyPub []byte
)

func init() {
	go func() {
		spaceKeyPri, spaceKeyPub, err := rsa.GenRsaKey(2048)
		if err != nil {
			fmt.Printf("\ninit GenRsaKey, spaceKey err:%v\n", err)
			logger.AppLogger().Debugf("init GenRsaKey, spaceKey err:%v", err)
		}
		preGeneratedSpaceKeyPri = spaceKeyPri
		preGeneratedSpaceKeyPub = spaceKeyPub

		passwordKeyPri, passwordKeyPub, err := rsa.GenRsaKey(2048)
		if err != nil {
			fmt.Printf("\ninit GenRsaKey, passwordKey err:%v\n", err)
			logger.AppLogger().Debugf("init GenRsaKey, passwordKey err:%v", err)
		}
		preGeneratedPasswordKeyPri = passwordKeyPri
		preGeneratedPasswordKeyPub = passwordKeyPub
	}()
}

func getSpaceKey(levelDBTrans *leveldb.Trans, aoId string) ([]byte, []byte, error) {
	logger.AppLogger().Debugf("getSpaceKey, aoId:%v",
		aoId)

	var err error
	exist, err := leveldb.Has(levelDBTrans, []byte(leveldb.KNameOfSpaceRSAPri(aoId)))
	if err != nil { // error
		logger.AppLogger().Debugf("getSpaceKey, leveldb.Has:%v, err:%v",
			leveldb.KNameOfSpaceRSAPri(aoId), err)
		return nil, nil, err
	}
	logger.AppLogger().Debugf("getSpaceKey, leveldb.Has:%v, exist:%v",
		leveldb.KNameOfSpaceRSAPri(aoId), exist)

	var priKeyBytes []byte
	var pubKeyBytes []byte
	if exist {
		priKeyBytes, err = leveldb.Get(levelDBTrans, []byte(leveldb.KNameOfSpaceRSAPri(aoId)))
		if err != nil { // error
			logger.AppLogger().Debugf("getSpaceKey, leveldb.Get:%v, err:%v",
				leveldb.KNameOfSpaceRSAPri(aoId), err)
			return nil, nil, err
		}

		pubKeyBytes, err = rsa.GetRsaPubKeyByPriKeyBytes(priKeyBytes)
		if err != nil {
			logger.AppLogger().Debugf("getPasswordKey, failed GetRsaPubKeyByPriKeyBytes, pubKeyBytes:%v, err:%v",
				pubKeyBytes, err)
			return nil, nil, err
		}

		logger.AppLogger().Debugf("getSpaceKey, GetRsaPubKeyByPriKeyBytes")

	} else { // not exist
		if preGeneratedSpaceKeyPri != nil && preGeneratedSpaceKeyPub != nil {
			priKeyBytes = preGeneratedSpaceKeyPri
			pubKeyBytes = preGeneratedSpaceKeyPub
			preGeneratedSpaceKeyPri = nil
			preGeneratedSpaceKeyPub = nil
			logger.AppLogger().Debugf("getSpaceKey, Using preGeneratedSpace")
		} else {
			priKeyBytes, pubKeyBytes, err = rsa.GenRsaKey(2048)
			if err != nil {
				return nil, nil, err
			}
			logger.AppLogger().Debugf("getSpaceKey, GenRsaKey : %v",
				leveldb.KNameOfSpaceRSAPri(aoId))
		}

		err = leveldb.Put(levelDBTrans, []byte(leveldb.KNameOfSpaceRSAPri(aoId)), priKeyBytes)
		if err != nil {
			return nil, nil, err
		}
	}

	return priKeyBytes, pubKeyBytes, nil
}

func aesDerivedKey(password string) []byte {
	logger.AppLogger().Debugf("aesDerivedKey")

	salt := sha256.Hash([]byte(saltSource), 1)[:32]
	dk := pbkdf2.Key([]byte(password), salt, 4096, 32, cryptosha256.New)
	return dk
}

func aesEncrypt(origData []byte, password string) ([]byte, error) {
	logger.AppLogger().Debugf("aesEncrypt")

	key := aesDerivedKey(password)
	iv := random.Random(passwordKeyIVLen)
	encbuf, err := aes.AesEncryptByPkcs5Padding(origData, key, iv)
	if err != nil {
		err1 := fmt.Errorf("AesEncryptByPkcs5Padding failed, err:%+v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return nil, err1
	}
	encData := make([]byte, passwordKeyVersionLen)
	encData = append(encData, iv...)
	encData = append(encData, encbuf...)
	return encData, nil
}

func aesDecrypt(origData []byte, password string) ([]byte, error) {
	logger.AppLogger().Debugf("aesDecrypt")

	key := aesDerivedKey(password)
	iv := origData[passwordKeyVersionLen : passwordKeyVersionLen+passwordKeyIVLen]
	encData := origData[passwordKeyVersionLen+passwordKeyIVLen:]
	decbuf, err := aes.AesDecryptByPkcs5Padding(encData, key, iv)
	if err != nil {
		err1 := fmt.Errorf("AesDecryptByPkcs5Padding failed, err:%+v", err)
		logger.AppLogger().Warnf("%+v", err1)
		return nil, err1
	}
	return decbuf, nil
}

func getEncryptedPrivatePasswordKey(levelDBTrans *leveldb.Trans, aoId string) ([]byte, bool, error) {
	logger.AppLogger().Debugf("getEncryptedPrivatePasswordKey, aoId:%v",
		aoId)

	var err error
	exist, err := leveldb.Has(levelDBTrans, []byte(leveldb.KNameOfPasswordRSAPri(aoId)))
	if err != nil { // error
		logger.AppLogger().Debugf("getPasswordKey, leveldb.Has:%v, err:%v",
			leveldb.KNameOfPasswordRSAPri(aoId), err)
		return nil, exist, err
	}
	fmt.Printf("\ngetPasswordKey, leveldb.Has:%v, err:%v\n",
		leveldb.KNameOfPasswordRSAPri(aoId), err)

	var encryptedPriKeyBytes []byte
	if exist {
		encryptedPriKeyBytes, err = leveldb.Get(levelDBTrans, []byte(leveldb.KNameOfPasswordRSAPri(aoId)))
		if err != nil { // error
			logger.AppLogger().Debugf("getPasswordKey, leveldb.Get:%v, err:%v",
				leveldb.KNameOfPasswordRSAPri(aoId), err)
			return nil, exist, err
		}
	}
	return encryptedPriKeyBytes, exist, nil
}

func getPasswordKey(levelDBTrans *leveldb.Trans, aoId, password string) ([]byte, []byte, []byte, error) {
	logger.AppLogger().Debugf("getPasswordKey, aoId:%v",
		aoId)

	var err error
	exist, err := leveldb.Has(levelDBTrans, []byte(leveldb.KNameOfPasswordRSAPri(aoId)))
	if err != nil { // error
		logger.AppLogger().Debugf("getPasswordKey, leveldb.Has:%v, err:%v",
			leveldb.KNameOfPasswordRSAPri(aoId), err)
		return nil, nil, nil, err
	}
	logger.AppLogger().Debugf("getPasswordKey, leveldb.Has:%v, exist:%v err:%v",
		leveldb.KNameOfPasswordRSAPri(aoId), exist, err)

	var priKeyBytes []byte
	var pubKeyBytes []byte
	var encryptedPriKeyBytes []byte
	if exist {
		encryptedPriKeyBytes, err = leveldb.Get(levelDBTrans, []byte(leveldb.KNameOfPasswordRSAPri(aoId)))
		if err != nil { // error
			logger.AppLogger().Debugf("getPasswordKey, leveldb.Get:%v, err:%v",
				leveldb.KNameOfPasswordRSAPri(aoId), err)
			return nil, nil, nil, err
		}
		logger.AppLogger().Debugf("getPasswordKey, len(encryptedPriKeyBytes):%v", encryptedPriKeyBytes)

		priKeyBytes, err := aesDecrypt(encryptedPriKeyBytes, password)
		if err != nil {
			return nil, nil, nil, err
		}
		logger.AppLogger().Debugf("getPasswordKey, aesDecrypt")

		pubKeyBytes, err = rsa.GetRsaPubKeyByPriKeyBytes(priKeyBytes)
		if err != nil {
			logger.AppLogger().Debugf("getPasswordKey, failed GetRsaPubKeyByPriKeyBytes, pubKeyBytes:%v, err:%v",
				pubKeyBytes, err)
			return nil, nil, nil, err
		}
		logger.AppLogger().Debugf("getPasswordKey, failed GetRsaPubKeyByPriKeyBytes, pubKeyBytes:%v",
			pubKeyBytes)

	} else { // not exist
		if preGeneratedPasswordKeyPri != nil && preGeneratedPasswordKeyPub != nil {
			priKeyBytes = preGeneratedPasswordKeyPri
			pubKeyBytes = preGeneratedPasswordKeyPub
			preGeneratedPasswordKeyPri = nil
			preGeneratedPasswordKeyPub = nil
			logger.AppLogger().Debugf("getPasswordKey, Using preGeneratedSpace")
		} else {
			priKeyBytes, pubKeyBytes, err = rsa.GenRsaKey(2048)
			if err != nil {
				return nil, nil, nil, err
			}
			logger.AppLogger().Debugf("getPasswordKey, GenRsaKey, pubKeyBytes:%v",
				leveldb.KNameOfPasswordRSAPri(aoId), string(pubKeyBytes))
		}

		encryptedPriKeyBytes, err = aesEncrypt(priKeyBytes, password)
		if err != nil {
			return nil, nil, nil, err
		}
		logger.AppLogger().Debugf("getPasswordKey, aesEncrypt")

		err := leveldb.Put(levelDBTrans, []byte(leveldb.KNameOfPasswordRSAPri(aoId)), encryptedPriKeyBytes)
		if err != nil {
			return nil, nil, nil, err
		}
		logger.AppLogger().Debugf("getPasswordKey, Put")
	}

	return encryptedPriKeyBytes, priKeyBytes, pubKeyBytes, nil
}

func deletePasswordKey(levelDBTrans *leveldb.Trans, aoId string) error {
	logger.AppLogger().Debugf("deletePasswordKey, aoId:%v",
		aoId)

	err := leveldb.Delete(levelDBTrans, []byte(leveldb.KNameOfPasswordRSAPri(aoId)))
	return err
}

func updatePasswordKey(levelDBTrans *leveldb.Trans, aoId string, oldPassword string, newPassword string) error {
	logger.AppLogger().Debugf("updatePasswordKey, aoId:%v",
		aoId)

	var err error
	exist, err := leveldb.Has(levelDBTrans, []byte(leveldb.KNameOfPasswordRSAPri(aoId)))
	if err != nil { // error
		return err
	}
	if !exist {
		return fmt.Errorf("key of %v not exist", leveldb.KNameOfPasswordRSAPri(aoId))
	}

	priKeyBytes, err := leveldb.Get(levelDBTrans, []byte(leveldb.KNameOfPasswordRSAPri(aoId)))
	if err != nil { // error
		logger.AppLogger().Debugf("updatePasswordKey, leveldb.Get:%v, err:%v",
			leveldb.KNameOfPasswordRSAPri(aoId), err)
		return err
	}

	decryptedPriKeyBytes, err := aesDecrypt(priKeyBytes, oldPassword)
	if err != nil {
		return err
	}

	pubKeyBytes, err := rsa.GetRsaPubKeyByPriKeyBytes(decryptedPriKeyBytes)
	if err != nil {
		logger.AppLogger().Debugf("updatePasswordKey, failed GetRsaPubKeyByPriKeyBytes, pubKeyBytes:%v, err:%v",
			pubKeyBytes, err)
		return err
	}
	logger.AppLogger().Debugf("updatePasswordKey, failed GetRsaPubKeyByPriKeyBytes, pubKeyBytes:%v",
		pubKeyBytes)

	encryptedPriKeyBytes, err := aesEncrypt(decryptedPriKeyBytes, newPassword)
	if err != nil {
		return err
	}

	err = leveldb.Put(levelDBTrans, []byte(leveldb.KNameOfPasswordRSAPri(aoId)), encryptedPriKeyBytes)
	if err != nil {
		return err
	}

	return err
}

func saveDidDoc(levelDBTrans *leveldb.Trans, did string, doc []byte) error {
	logger.AppLogger().Debugf("saveDidDoc, did:%v", did)

	if i := strings.Index(did, "?"); i > 0 {
		did = did[:i]
	}
	if i := strings.Index(did, "#"); i > 0 {
		did = did[:i]
	}
	logger.AppLogger().Debugf("saveDidDoc, did:%v",
		did)

	err := leveldb.Put(levelDBTrans, []byte(leveldb.KNameOfDidDoc(did)), doc)
	if err != nil {
		logger.AppLogger().Debugf("saveDidDoc, leveldb.Put:%v, err:%v",
			leveldb.KNameOfDidDoc(did), err)
		return err
	}
	return nil
}

func getDidDoc(levelDBTrans *leveldb.Trans, did string) ([]byte, error) {
	logger.AppLogger().Debugf("getDidDoc, did:%v", did)

	if i := strings.Index(did, "?"); i > 0 {
		did = did[:i]
	}
	if i := strings.Index(did, "#"); i > 0 {
		did = did[:i]
	}

	logger.AppLogger().Debugf("getDidDoc, did:%v",
		did)

	doc, err := leveldb.Get(levelDBTrans, []byte(leveldb.KNameOfDidDoc(did)))
	if err != nil {
		logger.AppLogger().Debugf("getDidDoc, leveldb.Put:%v, err:%v",
			leveldb.KNameOfDidDoc(did), err)
		return doc, err
	}
	return doc, nil
}
