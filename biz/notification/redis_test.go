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

package notification

import (
	"agent/config"
	"context"
	"fmt"
	"testing"

	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/go-redis/redis/v8"
)

func getRedisPortAndPassword() (string, string, error) {

	if fileutil.IsFileNotExist(config.Config.Box.RandDockercomposeRedisPort) {
		return "", "", fmt.Errorf("%v FileNotExist", config.Config.Box.RandDockercomposeRedisPort)
	}

	randstr, err := fileutil.ReadFromFile(config.Config.Box.RandDockercomposePassword)
	if err != nil {
		return "", "", fmt.Errorf("ReadFromFile %v err: %v", config.Config.Box.RandDockercomposePassword, err)
	}

	randRedisPort, err := fileutil.ReadFromFile(config.Config.Box.RandDockercomposeRedisPort)
	if err != nil {
		return "", "", fmt.Errorf("ReadFromFile %v err: %v", config.Config.Box.RandDockercomposeRedisPort, err)
	}

	port := string(randRedisPort)
	password := string(randstr)
	return port, password, nil
}

func TestStoreIntoRedis(t *testing.T) {
	port, password, err := getRedisPortAndPassword()
	if err != nil {
		t.Errorf("err:%v", err)
	}
	config.UpdateRedisConfig("127.0.0.1:"+port, password)

	clientUUID := "gotest1"
	optType := "upgrade_installing"
	id, err := storeIntoRedis(clientUUID, optType, "")
	if err != nil {
		t.Errorf("storeIntoRedis err:%v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       0, // use default DB
	})
	n, err := client.XDel(context.Background(), StreamNotification, id).Result()
	if err != nil {
		t.Errorf("err:%v", err)
	}
	if n != 1 {
		t.Errorf("n:%v NOT equal to 1", n)
	}
}
