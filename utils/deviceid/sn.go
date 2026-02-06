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

package deviceid

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

const (
	SnNumberLength    = 16 // SN 号长度
	ModelNumberLength = 8  // 设备内部代码长度, 需要用 "-"(减号) 或者 " "(空格) 或者 "0" 或者 "*"(星号) 来对字符串首或尾进行 padding。
)

// 设备代际号 | 扩展型号 | 年，十进制 | 月，十进制 | 流水号，十六进制 | 随机码，十六进制
// AS02 | 0 | 22 | 12 | 000001 | XXXX

type SnInfo struct {
	Sn             string `json:"sn"` // e.g.  2023010000010101
	DeviceModel    string `json:"deviceModel"`
	ProductionTime string `json:"productionTime"`
}

func GetSnNumber(snNumberStoreFile string) (string, error) {
	if fileutil.IsFileExist(snNumberStoreFile) {
		// fmt.Printf("GetSnNumber,%v exist\n", snNumberStoreFile)
		if content, err := fileutil.ReadFromFile(snNumberStoreFile); err != nil {
			fmt.Printf("GetSnNumber, ReadFromFile err: %v\n", err)
			return "", err
		} else {
			return string(content), nil
		}
	} else if strings.EqualFold(currentChip(), RPI) || strings.EqualFold(currentChip(), RK3568Dev) {
		fmt.Printf("GetSnNumber,%v not exist, so current hardware is chip:%v\n", snNumberStoreFile, currentChip())
		return "", fmt.Errorf("unsupported on this hardware platform")
	} else {
		fmt.Printf("GetSnNumber,%v not exist, so current hardware is rk's production\n", snNumberStoreFile)
		sn, err := getVendorSnNumber()
		if err != nil {
			fmt.Printf("GetSnNumber, getVendorSnNumber err: %v\n", err)
			return "", err
		} else {
			return sn, nil // 二代正式板
		}
	}
}

func getVendorSnNumber() (string, error) {
	data, err := VendorStorageRead(VENDOR_SN_ID)
	if err != nil {
		return "", fmt.Errorf("getVendorSnNumber, failed VendorStorageRead: %+v", err)
	}
	return string(data[:SnNumberLength]), nil
}

func GetModelNumber() (int, error) {
	data, err := VendorStorageRead(VENDOR_USER_NAME1)
	if err != nil {
		return -1, fmt.Errorf("GetModelNumber, failed VendorStorageRead: %+v", err)
	} else {
		s := string(data[:ModelNumberLength])
		s = strings.ReplaceAll(s, string([]byte{0x00}), "")
		s = strings.ReplaceAll(s, "-", "")
		s = strings.TrimSpace(s)
		s = strings.TrimLeft(s, "0")
		s = strings.ReplaceAll(s, "*", "")
		modelNumber, err := strconv.Atoi(s)
		if err != nil {
			return -2, err
		}

		return modelNumber, nil
	}
}
