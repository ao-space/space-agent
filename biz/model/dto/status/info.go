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
 * @Date: 2021-10-16 10:23:56
 * @LastEditors: wenchao
 * @LastEditTime: 2021-11-03 16:07:55
 * @Description:
 */
package status

import "agent/biz/model/device"

type Info struct {
	Status  string `json:"status"`  // 状态
	Version string `json:"version"` // 客户端版本

	IsClientPaired bool `json:"isClientPaired"` // 是否已经配对
	IsBoxInit      bool `json:"isBoxInit"`      // 是否已经初始化
	DockerStatus   int  `json:"dockerStatus"`   // Docker是否已经初始化完成

	TheBoxInfo        *device.DeviceInfo `json:"boxInfo"`
	TheBoxPriKeyBytes string             `json:"boxPriKeyBytes"`
	TheBoxPublicKey   string             `json:"boxPublicKey"`

	QrCode             string `json:"boxQrCode"`          // 绑定二维码
	TryoutCodeVerified bool   `json:"tryoutCodeVerified"` // 试用码是否验证通过(仅在 PC 试用场景下使用).
}
