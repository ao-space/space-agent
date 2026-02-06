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

package hardware

import (
	"fmt"
	"os"
	"strings"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

// Hardware        : BCM2835
// Hardware        : Firefly RK3568-ROC-PC HDMI (Linux)
func GetHardwareChip() (string, error) {
	hardware := ""
	b, err := fileutil.ReadFromFile("/proc/cpuinfo")
	if err != nil {
		err1 := fmt.Errorf("GetProcHardware, failed ReadFromFile, err:%v", err)
		return hardware, err1
	}
	s := string(b)
	// lines := tools.StringToLines(s)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")

	key := "Hardware"
	for _, line := range lines {
		// fmt.Printf("@@ line:%v, Contains:%v\n", line, strings.Contains(line, key))
		arr := strings.Split(line, ":")
		if len(arr) < 2 {
			continue
		}
		if strings.Contains(arr[0], key) {
			hardware = strings.TrimSpace(arr[1])
			// fmt.Printf("@@ hardware:%v\n", hardware)
			break
		}
	}
	return hardware, nil
}

// RunningInDocker is program running in container ?
func RunningInDocker() bool {
	// 由于其他方法在某些 OS 上有失效的可能性，暂时用用户传入的环境变量来判断。
	envkey := "AOSPACE_DATADIR"
	dataDir := os.Getenv(envkey)
	if len(dataDir) > 0 {
		fmt.Printf("RunningInDocker, dataDir:%v\n", dataDir)
		return true
	} else {
		fmt.Printf("RunningInDocker==false, dataDir:%v\n", dataDir)
		return false
	}
}
