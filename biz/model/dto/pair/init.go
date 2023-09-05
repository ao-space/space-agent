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
 * @Date: 2021-11-24 09:34:34
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-21 11:26:47
 * @Description:
 */

package pair

import "agent/biz/model/dto/device"

type InitResult struct {
	BoxName                string     `json:"boxName"`
	ClientUuid             string     `json:"clientUuid"`
	BoxUuid                string     `json:"boxUuid"`
	ProductId              string     `json:"productId"`
	Paired                 int        `json:"paired"`
	PairedBool             bool       `json:"pairedBool"`
	Connected              int        `json:"connected"`
	InitialEstimateTimeSec int        `json:"initialEstimateTimeSec"`
	Networks               []*Network `json:"network,omitempty"`
	SSPUrl                 string     `json:"sspUrl,omitempty"`
	NewBindProcessSupport  bool       `json:"newBindProcessSupport,omitempty"`

	device.BoxModelInfo
}
