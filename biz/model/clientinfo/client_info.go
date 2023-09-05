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

/*
 * @Author: wenchao
 * @Date: 2021-12-10 15:01:26
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-24 13:56:05
 * @Description:
 */

package clientinfo

import (
	"agent/config"
	"fmt"
	"strings"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/hash/sha256"
	"github.com/dungeonsnd/gocom/file/fileutil"
)

var clientPubKey string
var clientPriKey string
var sharedSecret string

type ClientKey struct {
	PubKey       []byte
	PrivKey      []byte
	SharedSecret string
}

func (ck *ClientKey) SaveKey() {
	err := fileutil.WriteToFile(config.Config.Box.ClientKey.RsaPubKeyFile, ck.PubKey, true)
	if err != nil {
		logger.AppLogger().Errorf("Write ClientKey.RsaPubKeyFile file failed, file:%v, err:%v",
			config.Config.Box.ClientKey.RsaPubKeyFile, err)
	} else {
		logger.AppLogger().Debugf("Write ClientKey.RsaPubKeyFile file succ, file:%v",
			config.Config.Box.ClientKey.RsaPubKeyFile)
	}
	err = fileutil.WriteToFile(config.Config.Box.ClientKey.RsaPriKeyFile, ck.PrivKey, true)
	if err != nil {
		logger.AppLogger().Errorf("Write ClientKey.RsaPriKeyFile file failed, file:%v, err:%v",
			config.Config.Box.ClientKey.RsaPriKeyFile, err)
	} else {
		logger.AppLogger().Debugf("Write ClientKey.RsaPriKeyFile file succ, file:%v",
			config.Config.Box.ClientKey.RsaPriKeyFile)
	}
}

func SetClientPubKey(theClientPubKey string) {
	clientPubKey = theClientPubKey
	logger.AppLogger().Debugf("SetClientPubKey, len(clientPubKey):%+v", len(clientPubKey))

	err := fileutil.WriteToFile(config.Config.Box.ClientKey.RsaPubKeyFile, []byte(clientPubKey), true)
	if err != nil {
		logger.AppLogger().Errorf("Write ClientKey.RsaPubKeyFile file failed, file:%v, err:%v",
			config.Config.Box.ClientKey.RsaPubKeyFile, err)
	} else {
		logger.AppLogger().Debugf("Write ClientKey.RsaPubKeyFile file succ, file:%v",
			config.Config.Box.ClientKey.RsaPubKeyFile)
	}
}

func SetClientPriKey(theClientPriKey string) {
	clientPriKey = theClientPriKey
	err := fileutil.WriteToFile(config.Config.Box.ClientKey.RsaPriKeyFile, []byte(clientPriKey), true)
	if err != nil {
		logger.AppLogger().Errorf("Write ClientKey.RsaPriKeyFile file failed, file:%v, err:%v",
			config.Config.Box.ClientKey.RsaPriKeyFile, err)
	} else {
		logger.AppLogger().Debugf("Write ClientKey.RsaPriKeyFile file succ, file:%v",
			config.Config.Box.ClientKey.RsaPriKeyFile)
	}
}

func GetClientPriKey() string {
	return clientPriKey
}

func SetSharedSecret(theSharedSecret string) error {
	sharedSecret = theSharedSecret

	err := fileutil.WriteToFile(config.Config.Box.ClientKey.SharedSecret, []byte(sharedSecret), true)
	if err != nil {
		logger.AppLogger().Errorf("Write ClientKey.SharedSecret file failed, file:%v, err:%v",
			config.Config.Box.ClientKey.SharedSecret, err)
		return fmt.Errorf("write SharedSecret failed, %v", err)
	} else {
		logger.AppLogger().Debugf("Write ClientKey.SharedSecret file succ, file:%v",
			config.Config.Box.ClientKey.SharedSecret)
		return nil
	}
}

func GetSharedSecret() (string, string, error) {
	if len(sharedSecret) < 1 {
		return "", "", fmt.Errorf("sharedSecret not setted")
	}

	sharedSecret = strings.ReplaceAll(sharedSecret, "\n", "")
	sharedSecret = strings.ReplaceAll(sharedSecret, "\r", "")
	sharedSecret = strings.ReplaceAll(sharedSecret, " ", "")
	rawIv := sha256.HashHex([]byte(sharedSecret), 1)[:16]
	return sharedSecret, rawIv, nil
}

func InitClientInfo() {
	logger.AppLogger().Debugf("InitClientInfo, GetAdminPairedStatus()= %v", GetAdminPairedStatus())
	if GetAdminPairedStatus() != DeviceAlreadyBound {
		return
	}

	if fileutil.IsFileExist(config.Config.Box.ClientKey.RsaPubKeyFile) {
		k, err := fileutil.ReadFromFile(config.Config.Box.ClientKey.RsaPubKeyFile)
		if err != nil {
			logger.AppLogger().Errorf("Read config.Config.Box.ClientKey.RsaPubKeyFile file failed, file:%v, err:%v",
				config.Config.Box.ClientKey.RsaPubKeyFile, err)
			return
		}
		logger.AppLogger().Debugf("Read config.Config.Box.ClientKey.RsaPubKeyFile file succ, file:%v, err:%v",
			config.Config.Box.ClientKey.RsaPubKeyFile, err)
		clientPubKey = string(k)
	} else {
		logger.AppLogger().Debugf("InitClientInfo, ClientKey.RsaPubKeyFile not exist, %v", config.Config.Box.ClientKey.RsaPubKeyFile)
	}

	if fileutil.IsFileExist(config.Config.Box.ClientKey.RsaPriKeyFile) {
		k, err := fileutil.ReadFromFile(config.Config.Box.ClientKey.RsaPriKeyFile)
		if err != nil {
			logger.AppLogger().Errorf("Read config.Config.Box.ClientKey.RsaPriKeyFile file failed, file:%v, err:%v",
				config.Config.Box.ClientKey.RsaPriKeyFile, err)
			return
		}
		logger.AppLogger().Debugf("Read config.Config.Box.ClientKey.RsaPriKeyFile file succ, file:%v, err:%v",
			config.Config.Box.ClientKey.RsaPriKeyFile, err)
		clientPriKey = string(k)
	} else {
		logger.AppLogger().Debugf("InitClientInfo, ClientKey.RsaPriKeyFile not exist, %v",
			config.Config.Box.ClientKey.RsaPriKeyFile)
	}

	if fileutil.IsFileExist(config.Config.Box.ClientKey.SharedSecret) {
		k, err := fileutil.ReadFromFile(config.Config.Box.ClientKey.SharedSecret)
		if err != nil {
			logger.AppLogger().Debugf("Read config.Config.Box.ClientKey.SharedSecret file failed, file:%v, err:%v",
				config.Config.Box.ClientKey.SharedSecret, err)
			return
		}
		logger.AppLogger().Debugf("Read config.Config.Box.ClientKey.SharedSecret file succ, file:%v, err:%v",
			config.Config.Box.ClientKey.SharedSecret, err)
		sharedSecret = string(k)
	} else {
		logger.AppLogger().Debugf("InitClientInfo, ClientKey.SharedSecret not exist, %v",
			config.Config.Box.ClientKey.SharedSecret)
	}

	logger.AppLogger().Infof("Leave InitClientInfo")
}
