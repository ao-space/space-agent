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
 * @Date: 2021-12-13 13:26:46
 * @LastEditors: jeffery
 * @LastEditTime: 2022-06-10 17:20:48
 * @Description:
 */

package clientinfo

import (
	"agent/config"
	"strconv"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

const (
	DeviceAlreadyBound = 0 // 已经配对
	NewDevice          = 1 // 新盒子
	DeviceRevoked      = 2 // 已解绑
)

type AdminPairedInfo struct {
	ClientUuid         string `json:"clientUUID"` // 客户端唯一id.
	AuthKey            string `json:"authKey"`
	ClientRegKey       string `json:"clientRegKey"`
	BoxName            string `json:"boxName"`
	UserDomain         string `json:"userDomain"`
	ClientPairedStatus string `json:"status"` // 0: 已经配对, 1: 新盒子, 2: 已解绑
}

func (a *AdminPairedInfo) AdminDomain() string {
	return a.UserDomain
}

func (a *AdminPairedInfo) Status() int {
	if a == nil || a.ClientPairedStatus == "" {
		return NewDevice
	}
	statusCode, err := strconv.Atoi(a.ClientPairedStatus)
	if err != nil {
		return NewDevice
	}
	return statusCode
}

func (a *AdminPairedInfo) Rebind() bool {
	return a.Status() == DeviceRevoked
}

func (a *AdminPairedInfo) HasPairedBefore() bool {
	return a.Status() == DeviceAlreadyBound || a.Status() == DeviceRevoked
}

func (a *AdminPairedInfo) AlreadyBound() bool {
	return a.Status() == DeviceAlreadyBound
}

func HasPairedBefore() bool {
	return GetAdminPairedStatus() == 0 || GetAdminPairedStatus() == 2
}

func GetAdminPairedStatus() int {
	pairedInfo := GetAdminPairedInfo()
	if pairedInfo == nil {
		logger.AppLogger().Debugf("GetAdminPairedStatus, GetAdminPairedInfo return nil")
		return NewDevice
	}

	// logger.AppLogger().Debugf("GetAdminPairedStatus, pairedInfo:%+v", pairedInfo)
	i, err := strconv.Atoi(pairedInfo.ClientPairedStatus)
	if err != nil {
		logger.AppLogger().Warnf("GetAdminPairedStatus, strconv.Atoi, pairedInfo.ClientPairedStatus, f=%v, err:%v",
			pairedInfo.ClientPairedStatus, err)
		return NewDevice
	}

	return i
}

func GetAdminPairedInfo() *AdminPairedInfo {
	f := config.Config.Box.BoxMetaAdminPair

	pairedInfo := &AdminPairedInfo{ClientUuid: "", ClientPairedStatus: "1"}
	// logger.AppLogger().Debugf("GetAdminPairedInfo, fileutil.IsFileExist(%v):%v", f, fileutil.IsFileExist(f))
	err := fileutil.ReadFileJsonToObject(f, pairedInfo)
	if err != nil {
		// logger.AppLogger().Debugf("ReadFileJsonToObject, f=%v, err:%v", f, err)
		return nil
	}
	// logger.AppLogger().Debugf("GetAdminPairedInfo, pairedInfo:%+v", pairedInfo)

	return pairedInfo
}

func GetAdminDomain() string {
	domain := ""

	pairedInfo := GetAdminPairedInfo()
	if pairedInfo == nil {
		return domain
	}
	domain = pairedInfo.UserDomain

	return domain
}
