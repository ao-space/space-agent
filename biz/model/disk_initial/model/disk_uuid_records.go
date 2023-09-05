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

package model

import (
	"agent/config"
	"strings"
	"time"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"github.com/dungeonsnd/gocom/file/fileutil"
)

const (
	RecordStartIdx = 1 // 记录号从多少开始
)

type DeviceUuidRecord struct {
	UpdatedTime string `json:"updatedTime"` // 更新时间

	Records []string `json:"records"` // 设备自增序列号所有记录, 包括曾经使用过的都会记录下来。
}

func NewDeviceUuidRecord() *DeviceUuidRecord {
	return &DeviceUuidRecord{UpdatedTime: time.Now().Format("2006-01-02 15:04:05"),
		Records: make([]string, 0)}
}

func (r *DeviceUuidRecord) Query(deviceUuid string) int64 {
	if r.Records != nil {
		for i, v := range r.Records {
			if strings.EqualFold(v, deviceUuid) {
				return int64(i + RecordStartIdx)
			}
		}
	}
	return -1
}

func (r *DeviceUuidRecord) CheckExistAndAppend(deviceUuid string) (int64, bool) {
	if r.Records != nil {
		for i, v := range r.Records {
			if strings.EqualFold(v, deviceUuid) {
				return int64(i + RecordStartIdx), true
			}
		}
	} else {
		r.Records = make([]string, 0)
	}

	r.Records = append(r.Records, deviceUuid)
	return int64(len(r.Records) - 1 + RecordStartIdx), false
}

func (r *DeviceUuidRecord) Write() error {

	f := config.Config.Box.Disk.DeviceUuidRecordFile

	r.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")

	jsonB, err := encoding.JsonEncode(r)
	if err != nil {
		logger.AppLogger().Errorf("DeviceUuidRecord, failed JsonEncode:%v, err:%v", r, err)
	} else {
		logger.AppLogger().Debugf("DeviceUuidRecord, json string: %+v", string(jsonB))
	}

	err = fileutil.WriteToFileAsJson(f, r, "  ", true)
	if err != nil {
		logger.AppLogger().Errorf("DeviceUuidRecord,  failed WriteToFileAsJson file:%v, err:%v", f, err)
		return err
	}
	logger.AppLogger().Debugf("DeviceUuidRecord WriteToFileAsJson succ, file:%v, info:%+v", f, r)

	return nil
}
