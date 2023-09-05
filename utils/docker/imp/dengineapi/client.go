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

package dengineapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var version = "1.39"

func SetClientVersion(v string) {
	version = v
}

func NewClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.WithVersion(version))
}

func Info() (string, error) {
	cli, err := NewClient()
	if err != nil {
		return "", fmt.Errorf("failed to call newClient, err:%v", err)
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to call cli.Info, err:%v", err)
	}

	return fmt.Sprintf("%+v", info), nil
}

// func Login(username, password, server string) error {
// 	cli, err := newClient()
// 	if err != nil {
// 		return fmt.Errorf("failed to call newClient, err:%v", err)
// 	}

// 	body, err := cli.RegistryLogin(context.Background(),
// 		types.AuthConfig{Username: username, Password: password, ServerAddress: server})

// 	fmt.Printf("RegistryLogin, body: %+v", body)
// 	return err
// }

func GenAuthString(username, password string) (string, error) {
	authConfig := types.AuthConfig{Username: username, Password: password}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	return authStr, nil
}
