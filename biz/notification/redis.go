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
 * @Author: jeffery
 * @Date: 2022-06-06 10:18:11
 * @LastEditors: jeffery
 * @LastEditTime: 2022-06-10 17:21:51
 * @Description:
 */

package notification

import (
	"agent/biz/model/clientinfo"
	"agent/config"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/random"
)

const (
	StreamNotification = "push_notification"
)

// clientUuid returns current admin client uuid
func clientUuid() (string, error) {

	pairedInfo := clientinfo.GetAdminPairedInfo()
	if pairedInfo == nil || pairedInfo.ClientPairedStatus != "0" {
		return "", fmt.Errorf("not paired, pairedInfo=%+v", pairedInfo)
	}
	clientUUID := pairedInfo.ClientUuid
	if len(clientUUID) < 1 {
		return clientUUID, fmt.Errorf("no ClientUuid to send push, pairedInfo=%+v", pairedInfo)
	}

	return clientUUID, nil
}

/*
"upgrade_downloaded_success" 下载成功 -- 管理员绑定端
"upgrade_installing" 正在安装 -- 除管理员绑定端以外所有端
"upgrade_success" 升级成功 -- 管理员绑定端
*/

// storeIntoRedis push msg to redis stream
func storeIntoRedis(clientUUID string, optType string, data interface{}) (string, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password, // no password set
		DB:       0,                            // use default DB
	})

	dataBytes, _ := json.Marshal(data)
	m := map[string]interface{}{"userId": "1",
		"clientUUID": clientUUID,
		"optType":    optType,
		"requestId":  random.GenUUID(),
		"data":       base64.StdEncoding.EncodeToString(dataBytes)}
	logger.NotificationLogger().Infof("storeIntoRedis,XAdd push_notification, map: %+v", m)
	id, err := client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: StreamNotification,
		Values: m,
	}).Result()

	if err != nil {
		logger.NotificationLogger().Warnf("Failed to send push_notification, err:%+v", err)
		return id, err
	} else {
		logger.NotificationLogger().Debugf("Succ to send push_notification, id:%+v", id)
		return id, nil
	}
}
