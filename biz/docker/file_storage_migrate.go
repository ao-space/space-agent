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

// 对于非二代板子的数据进行迁移

package docker

import (
	"agent/biz/model/device_ability"
	"agent/biz/model/disk_initial/model"
	"agent/config"
	"agent/utils/logger"
	"os"
	"time"

	"github.com/dungeonsnd/gocom/encrypt/random"
	"github.com/dungeonsnd/gocom/file/fileutil"
)

func MigrateFileStorageData() error {
	logger.DockerLogger().Debugf("MigrateFileStorageData")
	if fileutil.IsFileExist(config.Config.Box.Disk.DiskSharedInfoFile) {
		logger.DockerLogger().Debugf("MigrateFileStorageData, file exist %v", config.Config.Box.Disk.DiskSharedInfoFile)

		var v model.SharedDiskInfo
		err := fileutil.ReadFileJsonToObject(config.Config.Box.Disk.DiskSharedInfoFile, &v)
		if err == nil {
			logger.DockerLogger().Debugf("MigrateFileStorageData, ReadFileJsonToObject: %+v", v)

			if v.DiskMountInfos != nil && len(v.DiskMountInfos) > 0 {
				logger.DockerLogger().Debugf("MigrateFileStorageData, v.DiskMountInfos: %+v", v.DiskMountInfos)
				return nil
			}
		} else {
			logger.DockerLogger().Debugf("MigrateFileStorageData, ReadFileJsonToObject: %v , err.",
				config.Config.Box.Disk.DiskSharedInfoFile, err)
		}
	}

	if fileutil.IsFileExist(config.Config.Docker.ComposeFile) {
		err := DockerDownImmediately()
		if err != nil {
			logger.DockerLogger().Warnf("MigrateFileStorageData, failed DockerDownImmediately, err: %v , return.", err)
			return err
		}
	}
	logger.DockerLogger().Debugf("MigrateFileStorageData, succ DockerDownImmediately")

	srcDir := config.Config.Box.Disk.NoDisksFileStoragePath // /home/eulixspace/data
	if device_ability.GetAbilityModel().DeviceModelNumber <= device_ability.SN_GEN_CLOUD_DOCKER {
		srcDir = config.Config.Box.Disk.NoDisksFileStoragePathDockerDeploy // /home/eulixspace_link/data
	}
	logger.DockerLogger().Debugf("MigrateFileStorageData, srcDir:%v", srcDir)

	parts := config.Config.Box.Disk.StorageVolumePath // /home/eulixspace_file_storage/parts
	dstDir := fileutil.AddPathSepIfNeed(parts) + config.Config.Box.Disk.FileStorageVolumePathPrefix +
		config.Config.Box.Disk.StorageDummy // /home/eulixspace_file_storage/parts/bp_part_dummy

	// mv /home/eulixspace/data/dav /home/eulixspace_file_storage/parts/bp_part_dummy/
	// ...
	err := fileutil.CreateDirRecursive(dstDir)
	if err != nil {
		logger.DockerLogger().Warnf("MigrateFileStorageData, failed CreateDirRecursive (%v), err: %v , return.", dstDir, err)
		return err
	}
	logger.DockerLogger().Debugf("MigrateFileStorageData, succ CreateDirRecursive, parts:%v, dstDir:%v", parts, dstDir)

	srcFolderNames := []string{"dav", "eulixspace-files", "eulixspace-files-processed", "multipart", "third_party"}
	for _, srcFolderName := range srcFolderNames {
		srcFolder := fileutil.AddPathSepIfNeed(srcDir) + srcFolderName
		dstFolder := fileutil.AddPathSepIfNeed(dstDir) + srcFolderName
		logger.DockerLogger().Debugf("MigrateFileStorageData, srcFolder:%v, dstFolder:%v", srcFolder, dstFolder)
		if fileutil.IsFileExist(srcFolder) {
			err := os.Rename(srcFolder, dstFolder)
			if err != nil {
				logger.DockerLogger().Warnf("MigrateFileStorageData, failed Rename (%v >> %v), err: %v , return.", srcFolder, dstFolder, err)
				return err
			}
			logger.DockerLogger().Debugf("MigrateFileStorageData, succ Rename (%v >> %v), return.", srcFolder, dstFolder)
		}
	}

	return writeSharedDiskInfo()
}

func writeSharedDiskInfo() error {
	logger.DockerLogger().Debugf("writeSharedDiskInfo")

	// 更新磁盘初始化 共享文件
	tm := time.Now().Format("2006-01-02 15:04:05")
	dummyHwId := "0000a1d8"
	fShared := config.Config.Box.Disk.DiskSharedInfoFile
	sharedDiskInfo := &model.SharedDiskInfo{
		DiskInitialCode:             model.DiskInitialCode_Nomal,
		DiskInitialMessage:          "",
		DiskInitialProgress:         100,
		DiskExpandCode:              model.DiskExpandCode_NotExpanding,
		DiskExpandMessage:           "",
		DiskExpandProgress:          0,
		CreatedTime:                 tm,
		UpdatedTime:                 tm,
		RaidType:                    model.RaidType_Nomal,
		RaidDiskHwIds:               []string{},
		PrimaryStorageHwIds:         []string{dummyHwId},
		SecondaryStorageHwIds:       []string{},
		FileStorageVolumePathPrefix: config.Config.Box.Disk.FileStorageVolumePathPrefix,
	}

	sharedDiskInfo.DiskMountInfos = []*model.DiskMountInfo{{
		HwIds:                []string{dummyHwId},
		MountDevice:          "/dev/dummy",
		DeviceUuid:           random.GenUUID(),
		DeviceSequenceNumber: 1,
		MountPath:            "",
		DataFolderRoot:       config.Config.Box.Disk.StorageDummy,
		MapperName:           "",
		FSType:               "",
		IsPrimaryStorage:     true,
	}}
	err := fileutil.WriteToFileAsJson(fShared+".tmp", sharedDiskInfo, "  ", true)
	if err != nil {
		logger.DockerLogger().Errorf("writeSharedDiskInfo, fShared:%v, err:%v", fShared, err)
		return err
	}
	err = os.Rename(fShared+".tmp", fShared)
	if err != nil {
		logger.DockerLogger().Errorf("writeSharedDiskInfo, Rename (%v -> %v) err:%v", fShared+".tmp", fShared, err)
		return err
	}
	logger.DockerLogger().Debugf("writeSharedDiskInfo succ, fShared:%v, info:%+v", fShared, sharedDiskInfo)
	return nil
}
