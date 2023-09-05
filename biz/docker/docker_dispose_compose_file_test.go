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

package docker

import (
	"agent/config"
	"strings"
	"testing"

	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/dungeonsnd/gocom/log4go"
)

func TestWriteDefaultDockerComposeFile(t *testing.T) {
	writeDefaultDockerComposeFile()

	if fileutil.IsFileNotExist(config.Config.Box.RandDockercomposeRedisPort) {
		t.Errorf("%v FileNotExist", config.Config.Box.RandDockercomposeRedisPort)
	}

	randstr, err := fileutil.ReadFromFile(config.Config.Box.RandDockercomposePassword)
	if err != nil {
		log4go.E("ReadFromFile %v err: %v",
			config.Config.Box.RandDockercomposePassword, err)
		t.Errorf("ReadFromFile %v err: %v", config.Config.Box.RandDockercomposePassword, err)
	}

	randRedisPort, err := fileutil.ReadFromFile(config.Config.Box.RandDockercomposeRedisPort)
	if err != nil {
		log4go.E("ReadFromFile %v err: %v",
			config.Config.Box.RandDockercomposeRedisPort, err)
		t.Errorf("ReadFromFile %v err: %v", config.Config.Box.RandDockercomposeRedisPort, err)
	}

	target := "127.0.0.1:" + string(randRedisPort)
	if !strings.EqualFold(config.Config.Redis.Addr, target) {
		t.Errorf("Addr:%v NOT equal to %v", config.Config.Redis.Addr, target)
	}

	target = string(randstr)
	if !strings.EqualFold(config.Config.Redis.Password, target) {
		t.Errorf("Password:%v NOT equal to %v", config.Config.Redis.Password, target)
	}
}
