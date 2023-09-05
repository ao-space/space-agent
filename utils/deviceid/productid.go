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
 * @Date: 2021-11-10 14:41:02
 * @LastEditors: jeffery
 * @LastEditTime: 2022-05-18 15:46:30
 * @Description:
 */
package deviceid

import (
	hardware_util "agent/utils/hardware"
	"fmt"

	"github.com/dungeonsnd/gocom/encrypt/hash/sha256"
)

func GetProductId(cpuIdStoreFile string) (string, error) {
	if hardware_util.RunningInDocker() {
		cpuId, err := getCpuIdByRandom(cpuIdStoreFile)
		if err != nil {
			return "", err
		}
		s := fmt.Sprintf("eulixspace-productid-%v", cpuId)
		// fmt.Printf("GetProductId, RunningInDocker cpuId:%v\n", cpuId)
		return sha256.HashHex([]byte(s), 1)[:64], nil

	} else {
		cpuId, err := getCpuIdByCpuInfo(cpuIdStoreFile)
		if err != nil {
			return "", err
		}
		s := fmt.Sprintf("eulixspace-productid-%v", cpuId)
		// fmt.Printf("GetProductId, cpuId:%v\n", cpuId)
		return sha256.HashHex([]byte(s), 1), nil

	}
}
