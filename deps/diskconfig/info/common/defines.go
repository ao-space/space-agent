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

package common

type DiskInfo struct {
	HwId                  string   `json:"hwId"`                  // 生成的硬件id (格盘后不会变化). HASH(Model Family + Device Model +  Model Number + Serial Number)
	DeviceName            string   `json:"deviceName"`            // 设备名称, sda
	DevicePathName        string   `json:"devicePathName"`        // 设备名称, /dev/sda
	DefaultMapperName     string   `json:"defaultMapperName"`     // /dev/mapper 目录下的名称
	DefaultMapperPathName string   `json:"defaultMapperPathName"` // /dev/mapper 目录+名称
	MountName             string   `json:"mountName"`             // 挂载名称
	DisplayName           string   `json:"displayName"`           // UI上展示名称
	DiskUniId             string   `json:"diskUniId"`             // 磁盘唯一id (格盘后会变化)
	TransportType         int      `json:"transportType"`         // 传输类型 1: usb, 2: sata, 3: nvme
	ModelFamily           string   `json:"modelFamily"`           //
	DeviceModel           string   `json:"deviceModel"`           //
	SerialNumber          string   `json:"serialNumber"`          //
	ModelNumber           string   `json:"modelNumber"`           //
	BusNumber             int      `json:"busNumber"`             // BusNumber 硬盘总线号码. -1： unknown, 0: in sata 0; 1: in sata 4; 2: in sata 8; 101: m.2;
	PartedNames           []string `json:"partedNames"`           // 分区名称
	PartedUniIds          []string `json:"partedUniIds"`          // 分区唯一id (格盘后会变化)
}
