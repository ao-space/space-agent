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
	"agent/deps/diskconfig/info/common"
	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

type DiskInitialInfo struct {
	DiskConfig

	DiskInfos []*common.DiskInfo `json:"diskInfos"` // 初始化的每一块磁盘配置信息
}

const DefaultFsType = "ext4"

type DiskMountInfo struct {
	HwIds                []string `json:"hwIds"`               // e.g. "hwIds": ["52fa9e6c", "bd2c7570" ], 格式化之后不会变化
	MountDevice          string   `json:"mountDevice"`         // e.g. /dev/sda1, /dev/sdb1, /dev/vg0/lv0
	DeviceUuid           string   `json:"deviceUuid"`          // 格式化之后会变化
	DeviceSequenceNumber int64    `json:"dviceSequenceNumber"` // 设备自增序列号, 根据 DeviceUuid 判断，不重复时算新磁盘。
	MountPath            string   `json:"mountPath"`           // e.g. /mnt/bp/data/raid1_52fa9e6c-bd2c7570/mountpoint
	DataFolderRoot       string   `json:"dataFolderRoot"`      // raid1_138cf704-52fa9e6c , nvme_3841a657 in /mnt/bp/data/raid1_138cf704-52fa9e6c/mountpoint, /mnt/bp/data/nvme_3841a657/mountpoint"
	MapperName           string   `json:"mapperName"`          // e.g. luks_raid1_52fa9e6c-bd2c7570
	FSType               string   `json:"fSType"`              // ext4
	IsPrimaryStorage     bool     `json:"isPrimaryStorage"`    // 是否是主存储
}

type DiskConfig struct {
	DiskInitialCode     int    `json:"diskInitialCode"`     // 1: 磁盘正常; 2: 未初始化;3: 正在格式化; 4: 正在数据同步; 100: 未知错误; 101:  磁盘格式化错误; >101: 其他初始化错误;
	DiskInitialMessage  string `json:"diskInitialMessage"`  // 磁盘初始化结果/异常信息。
	DiskInitialProgress uint   `json:"diskInitialProgress"` // 磁盘初始化进度。

	DiskExpandCode     int      `json:"diskExpandCode"`     // 1: 扩容完成; 2: 未扩容状态; 3:正在扩容; 100: 扩容未知错误; 101:  扩容磁盘格式化错误; >101: 扩容其他错误;
	DiskExpandMessage  string   `json:"diskExpandMessage"`  // 磁盘扩容结果/异常信息。
	DiskExpandProgress uint     `json:"diskExpandProgress"` // 磁盘扩容进度。
	DiskExpandingHwIds []string `json:"diskExpandingHwIds"` // 正在扩容中的磁盘硬件 id 列表

	CreatedTime string `json:"createdTime"` // 创建时间
	UpdatedTime string `json:"updatedTime"` // 更新时间

	DiskEncrypt                int              `json:"diskEncrypt"`                // 磁盘加密与否。 1: 加密; 2: 不加密
	RaidType                   int              `json:"raidType"`                   // 1: normal; 2: raid1。
	RaidDiskHwIds              []string         `json:"raidDiskHwIds"`              // 参与 raid 的磁盘 id 列表.
	PrimaryStorageHwIds        []string         `json:"PrimaryStorageHwIds"`        // 主存储磁盘硬件 id 列表
	SecondaryStorageHwIds      []string         `json:"secondaryStorageHwIds"`      // 次存储磁盘硬件 id 列表
	PrimaryStorageMountPaths   string           `json:"primaryStorageMountPaths"`   // 主存储挂载目录
	SecondaryStorageMountPaths []string         `json:"secondaryStorageMountPaths"` // 次存储挂载目录
	MountInfos                 []*DiskMountInfo `json:"mountInfos"`                 // 挂载指令, 可用于重启后手动挂载
}

const (
	DiskInitialCode_Nomal                = 1   // 磁盘正常
	DiskInitialCode_NotInitialized       = 2   // 未初始化
	DiskInitialCode_Initializing         = 3   // 正在格式化
	DiskInitialCode_SynchronizingData    = 4   // 正在数据同步
	DiskInitialCode_UnkownError          = 100 // 初始化未知错误
	DiskInitialCode_FormatAndPartedError = 101 // 初始化磁盘格式化错误
	DiskInitialCode_DiskRaidError        = 102 // 初始化raid错误
	DiskInitialCode_DiskEncryptError     = 103 // 初始化加密错误
	DiskInitialCode_DiskMountError       = 104 // 初始化挂载错误
	DiskInitialCode_MigrateDockersError  = 105 // 初始化数据迁移错误

	DiskExpandCode_Nomal                      = 1   // 扩容完成
	DiskExpandCode_NotExpanding               = 2   // 未扩容状态
	DiskExpandCode_Expanding                  = 3   // 正在扩容
	DiskExpandCode_ExpandUnkownError          = 100 // 扩容未知错误
	DiskExpandCode_ExpandFormatAndPartedError = 101 // 扩容格式化错误
	DiskExpandCode_ExpandDiskRaidError        = 102 // 扩容raid错误
	DiskExpandCode_ExpandDiskEncryptError     = 103 // 扩容加密错误
	DiskExpandCode_ExpandDiskMountError       = 104 // 扩容挂载错误

)

const (
	RaidType_Nomal = 1 // normal
	RaidType_Raid1 = 2 // raid1
)

const (
	DiskEncrypt_Encrypt = 1 // 加密
	DiskEncrypt_Plain   = 2 // 不加密
)

var deviceUuidRecord *DeviceUuidRecord

var inprogressDiskInitialInfo *DiskInitialInfo

func init() {
	deviceUuidRecord = NewDeviceUuidRecord()
}

func SetInprogressDiskInitialInfo(info *DiskInitialInfo) {
	logger.AppLogger().Debugf("SetDiskInitialInfo, info: %+v", info)
	inprogressDiskInitialInfo = info
}

func ClearInprogressDiskInitialInfo() {
	logger.AppLogger().Debugf("ClearInprogressDiskInitialInfo")
	inprogressDiskInitialInfo = nil
}

func GetFileStoragePath() []string {
	var info *DiskInitialInfo
	if inprogressDiskInitialInfo != nil {
		info = inprogressDiskInitialInfo
		logger.AppLogger().Debugf("GetFileStoragePath, info: %+v", info)
	} else {
		info = ReadDiskInitialInfo()
		if info.DiskInitialCode != DiskInitialCode_Nomal {
			logger.AppLogger().Warnf("not initialized")
			return nil
		}
		logger.AppLogger().Debugf("GetFileStoragePath, ReadDiskInitialInfo info: %+v", info)
	}

	ret := make([]string, 0)
	for _, mountInfo := range info.MountInfos {
		target := mountInfo.MountPath + config.Config.Box.Disk.FileStorageInnerDataPath
		logger.AppLogger().Debugf("GetFileStoragePath, ReadDiskInitialInfo info: %+v", info)
		ret = append(ret, target)
	}
	return ret
}

func ReadDiskInitialInfo() *DiskInitialInfo {
	f := config.Config.Box.Disk.DiskInitialInfoFile

	info := &DiskInitialInfo{}
	info.DiskInitialCode = DiskInitialCode_NotInitialized
	info.DiskExpandCode = DiskExpandCode_NotExpanding
	err := fileutil.ReadFileJsonToObject(f, info)
	if err != nil {
		logger.AppLogger().Debugf("ReadFileJsonToObject, f=%v, err:%v", f, err)
	} else {
		// logger.AppLogger().Debugf("GetDiskInitialInfo, f=%v, info:%+v", f, info)
	}
	return info
}

//func WriteDiskInitialInfo(info *DiskInitialInfo) error {
//
//	// 更新 磁盘初始化文件
//	f := config.Config.Box.Disk.DiskInitialInfoFile
//
//	info.UpdatedTime = time.Now().Format("2006-01-02 15:04:05")
//
//	if fileutil.IsFileNotExist(f) {
//		info.CreatedTime = time.Now().Format("2006-01-02 15:04:05")
//	}
//
//	jsonB, err := encoding.JsonEncode(info)
//	if err != nil {
//		logger.AppLogger().Errorf("WriteDiskInitialInfo, failed JsonEncode:%v, err:%v", info, err)
//	} else {
//		logger.AppLogger().Debugf("WriteDiskInitialInfo, json string: %+v", string(jsonB))
//	}
//
//	err = fileutil.WriteToFileAsJson(f, info, "  ", true)
//	if err != nil {
//		logger.AppLogger().Errorf("WriteDiskInitialInfo, file:%v, err:%v", f, err)
//		return err
//	}
//	logger.AppLogger().Debugf("WriteDiskInitialInfo succ, file:%v, info:%+v", f, info)
//
//	return nil
//}
//
//func WriteDiskSharedInfoFile(info *DiskInitialInfo) error {
//
//	// 更新 设备自增序列号
//	updated := false
//	for _, mountInfos := range info.MountInfos {
//		deviceSequenceNumber, exist := deviceUuidRecord.CheckExistAndAppend(mountInfos.DeviceUuid)
//		logger.AppLogger().Debugf("WriteDiskInitialInfo, deviceSequenceNumber:%+v, exist:%+v, mountInfos:%+v",
//			deviceSequenceNumber, exist, mountInfos)
//		mountInfos.DeviceSequenceNumber = deviceSequenceNumber
//
//		if !exist {
//			updated = true
//		}
//	}
//	if updated {
//		deviceUuidRecord.Write()
//	}
//
//	// 更新磁盘初始化 共享文件
//	fShared := config.Config.Box.Disk.DiskSharedInfoFile
//	sharedDiskInfo := &SharedDiskInfo{
//		DiskInitialCode:             info.DiskInitialCode,
//		DiskInitialMessage:          info.DiskInitialMessage,
//		DiskInitialProgress:         info.DiskInitialProgress,
//		DiskExpandCode:              info.DiskExpandCode,
//		DiskExpandMessage:           info.DiskExpandMessage,
//		DiskExpandProgress:          info.DiskExpandProgress,
//		CreatedTime:                 info.CreatedTime,
//		UpdatedTime:                 info.UpdatedTime,
//		RaidType:                    info.RaidType,
//		RaidDiskHwIds:               info.RaidDiskHwIds,
//		PrimaryStorageHwIds:         info.PrimaryStorageHwIds,
//		SecondaryStorageHwIds:       info.SecondaryStorageHwIds,
//		FileStorageVolumePathPrefix: config.Config.Box.Disk.FileStorageVolumePathPrefix,
//	}
//	sharedDiskInfo.DiskMountInfos = info.MountInfos
//	err := fileutil.WriteToFileAsJson(fShared, sharedDiskInfo, "  ", true)
//	if err != nil {
//		logger.AppLogger().Errorf("WriteDiskInitialInfo, fShared:%v, err:%v", fShared, err)
//		return err
//	}
//	logger.AppLogger().Debugf("WriteDiskInitialInfo succ, fShared:%v, info:%+v", fShared, sharedDiskInfo)
//
//	return nil
//}
