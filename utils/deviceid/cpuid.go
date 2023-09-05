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
	"regexp"
	"strings"

	"agent/utils/logger"
	"github.com/dungeonsnd/gocom/encrypt/hash/sha256"
	"github.com/dungeonsnd/gocom/encrypt/random"
	"github.com/dungeonsnd/gocom/file/fileutil"
)

func getCpuIdByRandom(cpuIdStoreFile string) (string, error) {
	if fileutil.IsFileExist(cpuIdStoreFile) {
		if content, err := fileutil.ReadFromFile(cpuIdStoreFile); err != nil {
			return "", err
		} else {
			return string(content), nil
		}
	}

	r := sha256.HashHex([]byte(random.GenUUID()), 1)[:64]
	return r, fileutil.WriteToFile(cpuIdStoreFile, []byte(r), false)
}

func getCpuIdByCpuInfo(cpuIdStoreFile string) (string, error) {

	if fileutil.IsFileExist(cpuIdStoreFile) {
		if content, err := fileutil.ReadFromFile(cpuIdStoreFile); err != nil {
			return "", nil
		} else {
			return string(content), nil
		}
	}

	b, err := fileutil.ReadFromFile("/proc/cpuinfo")
	if err != nil {
		logger.AppLogger().Warnf("GetBtId, failed ReadFromFile, err:%v", err)
		return "", fmt.Errorf("failed ReadFromFile, err:%v", err)
	}
	s := string(b)
	r, err := regexp.Compile(".*Serial.*[a-zA-Z0-9]*")
	if err != nil {
		return "", fmt.Errorf("failed regexp.Compile r, err:%v", err)
	}
	ser := r.FindString(s)
	if len(ser) > 1 {
		ser = strings.Replace(ser, "Serial", "", -1)
		ser = strings.Replace(ser, ":", "", -1)
		ser = strings.Replace(ser, "\t", "", -1)
		ser = strings.Replace(ser, " ", "", -1)
		ser = strings.Replace(ser, "\r", "", -1)
		ser = strings.Replace(ser, "\n", "", -1)
		return ser, fileutil.WriteToFile(cpuIdStoreFile, []byte(ser), false)
	}
	return "", fmt.Errorf("not found Serial")
}
