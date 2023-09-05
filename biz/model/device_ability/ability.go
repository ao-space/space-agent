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
	"agent/utils/logger"
)

const (
	SN_GEN_1                       = 100 // 第一代
	SN_GEN_2                       = 200 // 第二代
	SN_SUPPORTED_FROM_MODEL_NUMBER = 210 // 从哪个型号开始支持 sn 号

	SN_GEN_VM           = -100 // 虚拟机版本
	SN_GEN_CLOUD_DOCKER = -200 // 云试用容器版本
	SN_GEN_PC_DOCKER    = -300 // PC容器版本
)

// DeviceAbility 设备支持能力模型
type DeviceAbility struct {
	DeviceModelNumber            int    `json:"deviceModelNumber"` // 产品型号数字(内部使用, 1xx: 树莓派, 2xx: 二代, ...)
	SnNumber                     string `json:"snNumber"`
	InnerDiskSupport             bool   `json:"innerDiskSupport"`    // 内部磁盘支持(SATA 和 m.2)
	SupportUSBDisk               bool   `json:"supportUSBDisk"`      // 是否支持 USB 磁盘. 210型号 及以上不支持
	SecurityChipSupport          bool   `json:"securityChipSupport"` // 支持加密芯片
	RunInDocker                  bool   `json:"runInDocker"`         // 以容器方式运行
	OpenSource                   bool   `json:"openSource"`
	BluetoothSupport             bool   `json:"bluetoothSupport"`             // 当前设备是否支持蓝牙
	NetworkConfigSupport         bool   `json:"networkConfigSupport"`         // 当前设备是否支持网络配置
	LedSupport                   bool   `json:"ledSupport"`                   // 当前设备是否支持 Led
	BackupRestoreSupport         bool   `json:"backupRestoreSupport"`         // 当前设备是否支持备份恢复
	AospaceAppSupport            bool   `json:"aospaceappSupport"`            // 当前设备是否支持傲空间应用 PC: true
	AospaceDevOptionSupport      bool   `json:"aospaceDevOptionSupport"`      // 是否支持开发者选项 PC: true
	AospaceSwitchPlatformSupport bool   `json:"aospaceSwitchPlatformSupport"` // 是否支持切换平台 PC: false
	UpgradeApiSupport            bool   `json:"upgradeApiSupport"`            // 当前设备是否支持升级API
}

// 设备支持能力
var deviceAbility *DeviceAbility

func init() {
	deviceAbility = &DeviceAbility{}

	deviceAbility.DeviceModelNumber = getDeviceModelNumber()

	// if config.Config.EnableSecurityChip {
	deviceAbility.SecurityChipSupport = supportSecurityChip()
	// } else {
	// 	deviceAbility.SecurityChipSupport = false
	// }

	snNumber, err := deviceid.GetSnNumber(config.Config.Box.SnNumberStoreFile)
	if err != nil {
		logger.AppLogger().Debugf("failed GetSnNumber, err:%v", err)
	}
	deviceAbility.SnNumber = snNumber

	deviceAbility.InnerDiskSupport = supportInnerDisk()
	deviceAbility.SupportUSBDisk = supportUSBDisk()

	deviceAbility.RunInDocker = false
	if getDeviceModelNumber() <= SN_GEN_CLOUD_DOCKER {
		deviceAbility.RunInDocker = true
	}

	deviceAbility.BluetoothSupport = true
	if getDeviceModelNumber() <= SN_GEN_VM {
		deviceAbility.BluetoothSupport = false
	}

	deviceAbility.NetworkConfigSupport = true
	if getDeviceModelNumber() <= SN_GEN_VM {
		deviceAbility.NetworkConfigSupport = false
	}

	deviceAbility.LedSupport = false
	if getDeviceModelNumber() >= SN_SUPPORTED_FROM_MODEL_NUMBER {
		deviceAbility.LedSupport = true
	}

	deviceAbility.BackupRestoreSupport = true
	if getDeviceModelNumber() <= SN_GEN_CLOUD_DOCKER {
		deviceAbility.BackupRestoreSupport = config.Config.EnableBackupRestoreSupportWhenRunAsDocker
	}

	deviceAbility.AospaceAppSupport = true
	deviceAbility.AospaceDevOptionSupport = true
	deviceAbility.AospaceSwitchPlatformSupport = true
	deviceAbility.OpenSource = true
	deviceAbility.UpgradeApiSupport = true
}

func GetAbilityModel() *DeviceAbility {
	// fmt.Printf("@@ GetAbilityModel: %+v \n", deviceAbility)
	snNumber, err := deviceid.GetSnNumber(config.Config.Box.SnNumberStoreFile)
	if err != nil {
		// logger.AppLogger().Debugf("failed GetSnNumber, err:%v", err)
	}
	deviceAbility.SnNumber = snNumber
	return deviceAbility
}
