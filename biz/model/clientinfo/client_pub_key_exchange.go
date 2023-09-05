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
 * @Date: 2021-12-13 14:03:26
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-24 13:54:47
 * @Description:
 */

package clientinfo

import "agent/utils/logger"

var clientExchangePubKey string
var clientExchangePriKey string

func SetClientExchangePubKey(theClientPubKey string) {
	clientExchangePubKey = theClientPubKey
	logger.AppLogger().Debugf("SetClientExchangePubKey, clientExchangePubKey: %+v", clientExchangePubKey)
}

func GetClientExchangePubKey() string {
	logger.AppLogger().Debugf("GetClientExchangePubKey, clientExchangePubKey: %+v", clientExchangePubKey)
	return clientExchangePubKey
}

func SetClientExchangePriKey(theClientPriKey string) {
	clientExchangePriKey = theClientPriKey
}

func ClientExchangePubKeyExchanged() bool {
	logger.AppLogger().Debugf("ClientExchangePubKeyExchanged, clientExchangePubKey: %+v", clientExchangePubKey)
	return len(clientExchangePubKey) > 0
}

func SaveClientExchangeKey() {
	SetClientPubKey(clientExchangePubKey)
	SetClientPriKey(clientExchangePriKey)
}
