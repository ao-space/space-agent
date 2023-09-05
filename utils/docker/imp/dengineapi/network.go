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
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func CreateNetwork(cli *client.Client, name string) (string, string, error) {
	var err error
	if cli == nil {
		cli, err = NewClient()
		if err != nil {
			return "", "", fmt.Errorf("failed NewClient, err:%v", err)
		}
		defer cli.Close()
	}

	createOpts := types.NetworkCreate{
		CheckDuplicate: true,
		Internal:       false,
	}
	retNetwork, err := cli.NetworkCreate(context.Background(), name, createOpts)
	if err != nil {
		return "", "", err
	}

	// fmt.Printf("retNetwork: %+v \n", retNetwork.ID, retNetwork.Warning)

	return retNetwork.ID, retNetwork.Warning, nil
}
