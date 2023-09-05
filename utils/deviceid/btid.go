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
 * @LastEditTime: 2022-04-14 17:00:56
 * @Description:
 */
package deviceid

import (
	"agent/config"
	hardware_util "agent/utils/hardware"
	"fmt"
	"os"
	"strings"

	"github.com/dungeonsnd/gocom/encrypt/hash/sha256"
)

const (
	RPI       = "BCM2835"
	RK3568Dev = "Firefly RK3568-ROC-PC HDMI (Linux)"
)

type BtidContainer struct {
	SnNumberStoreFile string
}

type BtidDevice struct {
	CpuIdFile string
}

type BtidService interface {
	Get() (string, error)
	Hash(sn string) string
}

func (bc *BtidContainer) Get() (string, error) {
	sn, err := GetSnNumber(bc.SnNumberStoreFile)
	if err != nil {
		return "", err
	}
	return HashHex(sn), nil // pc docker
}

func (bd *BtidDevice) Get() (string, error) {
	var (
		sn   string
		err1 error
	)
	if _, err := os.Stat(VendorStorageFile); err != nil {
		sn, err1 = getCpuIdByCpuInfo(bd.CpuIdFile)
		if err1 != nil {
			return "", err1
		}
		return HashHex(sn), err1
	}
	sn, err1 = getVendorSnNumber()
	if err1 != nil {
		return "", err1
	}
	return HashHex(sn), nil
}

func (bc *BtidContainer) Hash(sn string) string {
	return HashHex(sn + config.Config.Box.Avahi.BtIdHashPrefix)
}

func (bd *BtidDevice) Hash(sn string) string {
	return HashHex(sn + config.Config.Box.Avahi.BtIdHashPrefix)
}

func HashHex(sn string) string {
	return sha256.HashHex([]byte(sn), 1)[:16]
}

func GetBtId(cpuIdStoreFile string, snNumberStoreFile string) (string, error) {
	if hardware_util.RunningInDocker() {
		sn, err := GetSnNumber(snNumberStoreFile)
		if err != nil {
			return "", err
		} else {
			return sha256.HashHex([]byte(sn), 1)[:16], nil // pc docker
		}

	} else if strings.EqualFold(currentChip(), RPI) || strings.EqualFold(currentChip(), RK3568Dev) {
		cpuId, err := getCpuIdByCpuInfo(cpuIdStoreFile)
		if err != nil {
			return "", err
		}
		s := fmt.Sprintf("eulixspace-btid-%v", cpuId)
		return sha256.HashHex([]byte(s), 1)[:16], nil

	} else {
		sn, err := getVendorSnNumber()
		if err != nil {
			return "", err
		} else {
			return sha256.HashHex([]byte(sn), 1)[:16], nil // 二代正式板
		}
	}
}
