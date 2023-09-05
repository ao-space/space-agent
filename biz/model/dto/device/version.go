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

package device

import (
	"agent/biz/model/device_ability"
	"agent/utils/docker/dockermodel"
)

type ServiceVersion struct {
	Created     int64  `json:"created"`
	ServiceName string `json:"serviceName"`
	Version     string `json:"version"`
}

type BoxDeviceVersion struct {
	BoxModelInfo
	SnNumber       string            `json:"snNumber"`       // sn号
	ServiceVersion []*ServiceVersion `json:"serviceVersion"` // 容器版本
	// ServiceDetail  []*dockermodel.DockerImage `json:"serviceDetail"`  // 容器版本
	ServiceDetail []*dockermodel.DockerContainer `json:"serviceDetail"` // 容器版本
}

// 中文名称。
func GetDeviceName() string {
	if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_GEN_2 {
		return "傲空间（第二代）"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_VM {
		return "傲空间（PC版）"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_CLOUD_DOCKER {
		return "傲空间（在线版）"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_PC_DOCKER {
		return "傲空间（PC版）"
	} else {
		return "傲空间（第一代）"
	}
}

// 英文名称。
func GetDeviceNameEn() string {
	if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_GEN_2 {
		return "AS (2nd Generation)"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_VM {
		return "AS (PC)"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_CLOUD_DOCKER {
		return "AS (Online)"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_PC_DOCKER {
		return "AS (PC)"
	} else {
		return "AS (1nd Generation)"
	}
}

// 英文代系
// (1nd Generation)
func GetGenerationEn() string {
	if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_GEN_2 {
		return "(2nd Generation)"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_VM {
		return "(PC)"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_CLOUD_DOCKER {
		return "(Online)"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_PC_DOCKER {
		return "(PC)"
	} else {
		return "(1nd Generation)"
	}
}

// 中文代系
// （第一代）
func GetGenerationZh() string {
	if device_ability.GetAbilityModel().DeviceModelNumber >= device_ability.SN_GEN_2 {
		return "（第二代）"
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_VM {
		return ""
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_CLOUD_DOCKER {
		return ""
	} else if device_ability.GetAbilityModel().DeviceModelNumber == device_ability.SN_GEN_PC_DOCKER {
		return ""
	} else {
		return "（第一代）"
	}
}

// 产品型号。一代为空，二代才有。
func GetProductModel() string {
	return ""
}

// sn 号
func GetSnNumber() string {
	return device_ability.GetAbilityModel().SnNumber
}
