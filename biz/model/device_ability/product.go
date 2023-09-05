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

package device_ability

import (
	"agent/config"
	"agent/utils/deviceid"
	hardware_util "agent/utils/hardware"
	"strings"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

// cat /proc/cpuinfo
//
// 树莓派:
// Hardware        : BCM2835
//
// 开发板:
// Hardware        : Firefly RK3568-ROC-PC HDMI (Linux)
//
// 正式板
// Hardware        : Rockchip RK3568 EVB1 DDR4 V10 Board
//

const (
	RPI           = "BCM2835"
	RK3568        = "Firefly RK3568-ROC-PC HDMI (Linux)"
	RK3568Product = "Rockchip RK3568 EVB1 DDR4 V10 Board"
)

func currentChip() string {
	rt, err := hardware_util.GetHardwareChip()
	if err != nil {
		logger.AppLogger().Warnf("failed GetHardwareChip, err:%v", err)
		return ""
	}
	return rt
}

func getDeviceModelNumber() int {
	if hardware_util.RunningInDocker() { // PC 容器版本没有型号存储区, 判断是不是运行在容器中。
		// fmt.Printf("getDeviceModelNumber RunningInDocker \n")
		if fileutil.IsFileExist(config.Config.Box.SnNumberStoreFile) {
			b, err := fileutil.ReadFromFile(config.Config.Box.SnNumberStoreFile)
			if len(string(b)) > config.Config.Box.SnNumberModelLength {
				if err == nil && strings.Index(string(b)[:config.Config.Box.SnNumberModelLength], config.Config.Box.SnNumberModelContent) == 0 {
					return SN_GEN_CLOUD_DOCKER
				}
			}
		}

		return SN_GEN_PC_DOCKER
	} else if strings.EqualFold(currentChip(), RPI) { // 树莓派没有型号存储区, 检测 芯片类型 来判断。
		// fmt.Printf("getDeviceModelNumber RPI \n")
		return SN_GEN_1
	} else if strings.EqualFold(currentChip(), RK3568) { // 二代开发板没有型号存储区, 检测 芯片类型 来判断。
		// fmt.Printf("getDeviceModelNumber RK3568 \n")
		return SN_GEN_2
	}

	// fmt.Printf("getDeviceModelNumber deviceid.GetModelNumber \n")
	modelNumber, err := deviceid.GetModelNumber() // 二代正式板有型号存储区，直接读取型号代码。
	if err != nil {
		logger.AppLogger().Warnf("failed GetModelNumber, err:%v", err)
		return 0
	} else {
		return modelNumber
	}

}

// 是否支持磁盘
func supportInnerDisk() bool {
	if hardware_util.RunningInDocker() {
		return false
	} else if strings.EqualFold(currentChip(), RPI) {
		return false
	}
	return true
}

// 是否支持USB磁盘
func supportUSBDisk() bool {
	if getDeviceModelNumber() >= SN_GEN_2 && getDeviceModelNumber() < SN_SUPPORTED_FROM_MODEL_NUMBER {
		return true
	}
	return false
}

// 是否支持加密芯片
func supportSecurityChip() bool {
	if getDeviceModelNumber() >= SN_GEN_2 {
		return true
	}
	return false
}
