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
	"agent/config"
)

const (
	FirstGenDevNameEn        = "AS (1st Generation)"
	SecondGenDevNameEn       = "AS (2nd Generation)"
	PCVersionDevNameEn       = "AS (PC)"
	OnlineVersionDevNameEn   = "AS (Online)"
	FirstGenDevNameZhCn      = "傲空间（第一代)"
	SecondGenDevNameZhCn     = "傲空间（第二代)"
	PCVersionDevNameZhCn     = "傲空间（PC)"
	OnlineVersionDevNameZhCn = "傲空间（在线版)"
	FirstGenEn               = "(1st Generation)"
	SecondGenEn              = "(2nd Generation)"
	PCVersionEn              = "(PC)"
	OnlineVersionEn          = "(Online)"
	FirstGenZhCn             = "（第一代)"
	SecondGenZhCn            = "（第二代)"
	PCVersionZhCn            = "（PC)"
	OnlineVersionZhCn        = "(在线版)"
)

type BoxModelInfo struct {
	DeviceName   string `json:"deviceName"`   // 设备名称
	DeviceNameEn string `json:"deviceNameEn"` // 英文名称

	GenerationEn string `json:"generationEn"` // 英文代系
	GenerationZh string `json:"generationZh"` // 中文代系

	ProductModel string `json:"productModel"` // 产品型号

	SpaceVersion string `json:"spaceVersion"` // 傲空间版本
	OSVersion    string `json:"osVersion"`    // 操作系统版本

	DeviceLogoUrl string `json:"deviceLogoUrl"` // 设备图片链接

	DeviceAbility *device_ability.DeviceAbility `json:"deviceAbility"` // 设备能力模型
	//OpenSource    bool                          `json:"openSource"`
}

var boxModelMap = make(map[int]*BoxModelInfo)

func init() {
	boxModelMap[device_ability.SN_GEN_2] = &BoxModelInfo{
		DeviceName:    SecondGenDevNameZhCn,
		DeviceNameEn:  SecondGenDevNameEn,
		GenerationEn:  SecondGenEn,
		GenerationZh:  SecondGenZhCn,
		ProductModel:  "",
		SpaceVersion:  "",
		OSVersion:     "",
		DeviceLogoUrl: "",
		DeviceAbility: nil,
	}

	boxModelMap[device_ability.SN_GEN_1] = &BoxModelInfo{
		DeviceName:    FirstGenDevNameZhCn,
		DeviceNameEn:  FirstGenDevNameEn,
		GenerationEn:  FirstGenEn,
		GenerationZh:  FirstGenZhCn,
		ProductModel:  "",
		SpaceVersion:  "",
		OSVersion:     "",
		DeviceLogoUrl: "",
		DeviceAbility: nil,
	}

	boxModelMap[device_ability.SN_GEN_PC_DOCKER] = &BoxModelInfo{
		DeviceName:    PCVersionDevNameZhCn,
		DeviceNameEn:  PCVersionDevNameEn,
		GenerationEn:  PCVersionEn,
		GenerationZh:  PCVersionZhCn,
		ProductModel:  "",
		SpaceVersion:  "",
		OSVersion:     "",
		DeviceLogoUrl: "",
		DeviceAbility: nil,
	}

	boxModelMap[device_ability.SN_GEN_VM] = boxModelMap[device_ability.SN_GEN_PC_DOCKER]

	boxModelMap[device_ability.SN_GEN_CLOUD_DOCKER] = &BoxModelInfo{
		DeviceName:    OnlineVersionDevNameZhCn,
		DeviceNameEn:  OnlineVersionDevNameEn,
		GenerationEn:  OnlineVersionEn,
		GenerationZh:  OnlineVersionZhCn,
		ProductModel:  "",
		SpaceVersion:  config.VersionNumber,
		OSVersion:     "",
		DeviceLogoUrl: "",
		DeviceAbility: nil,
	}
}
