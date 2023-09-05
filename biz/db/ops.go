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

package db

import (
	"agent/biz/model/upgrade"
	"agent/utils/logger"
	"fmt"
	"github.com/imdario/mergo"
	"time"
)

// ReadTask is to read a task from the database, when versionId is set, It will match the version.
func ReadTask(versionId string) (*upgrade.Task, error) {
	defer lock.RUnlock()
	lock.RLock()
	task := new(upgrade.Task)
	db, err := NewDBClient()
	if err != nil {
		return task, err
	}
	err = db.Read(conf.UpgradeCollection, conf.TaskResource, &task)
	if err != nil {
		return task, err
	}
	if versionId != "" && versionId != task.VersionId {
		return task, fmt.Errorf("no record of the specified version exists")
	}
	return task, nil
}

func MarkTaskInstalling(versionId string) (*upgrade.Task, error) {
	logger.UpgradeLogger().Debugf("Marking task %s status %s", versionId, upgrade.Installing)
	doc, err := UpdateOrCreateTask(&upgrade.Task{
		VersionId:        versionId,
		Status:           upgrade.Installing,
		StartInstallTime: time.Now().Format(time.RFC3339)})
	return doc, err
}

func MarkTaskRebootFlag(versionId string) (*upgrade.Task, error) {
	logger.UpgradeLogger().Debugf("Marking task %s status %s", versionId, upgrade.Installing)
	doc, err := UpdateOrCreateTask(&upgrade.Task{
		NeedReboot: true})
	return doc, err
}

// UpdateOrCreateTask to write db for create a new task or update old task, according to versionId.
func UpdateOrCreateTask(newT *upgrade.Task) (*upgrade.Task, error) {
	defer lock.Unlock()
	lock.Lock()
	task := new(upgrade.Task)
	db, err := NewDBClient()
	if err != nil {
		return task, err
	}
	err = db.Read(conf.UpgradeCollection, conf.TaskResource, &task)
	if err != nil {
		return task, err
	}

	if newT.VersionId != task.VersionId {
		// 一个新的 Task， 需要覆盖空值
		err = db.Write(conf.UpgradeCollection, conf.TaskResource, newT)
		if err != nil {
			return task, fmt.Errorf("update task => %w", err)
		}
		return newT, nil

	} else {
		// 相同的 task，空值不覆盖
		err = mergo.Merge(task, newT, mergo.WithTransformers(upgrade.TimeTransformer{}), mergo.WithOverride)
		if err != nil {
			return task, fmt.Errorf("update task => %w", err)
		}
		err = db.Write(conf.UpgradeCollection, conf.TaskResource, task)
		if err != nil {
			return task, fmt.Errorf("update task => %w", err)
		}
		return task, nil
	}
}

func MarkTaskDownErr(versionId string) *upgrade.Task {
	logger.UpgradeLogger().Debugf("Marking task %s status %s", versionId, upgrade.DownloadErr)
	doc, err := UpdateOrCreateTask(&upgrade.Task{
		VersionId:       versionId,
		Status:          upgrade.DownloadErr,
		InstallStatus:   upgrade.Err,
		DoneInstallTime: time.Now().Format(time.RFC3339)})
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to mark task download err %s", err)
	}
	return doc
}

func MarkTaskDownloaded(versionId string, rpmInfo upgrade.VersionDownInfo, imageInfo upgrade.VersionDownInfo, kernelInfo upgrade.VersionDownInfo) *upgrade.Task {
	logger.UpgradeLogger().Debugf("Marking task %s status %s", versionId, upgrade.Downloaded)
	task := upgrade.Task{
		VersionId:    versionId,
		RpmPkg:       rpmInfo,
		ContainerImg: imageInfo,
		KernelImg:    kernelInfo,
		Status:       upgrade.Downloaded,
		DownStatus:   upgrade.Done,
		DoneDownTime: time.Now().Format(time.RFC3339),
	}
	doc, err := UpdateOrCreateTask(&task)

	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to mark task downloaded %s", err)
	}
	return doc
}

func MarkTaskInstallErr(versionId string) *upgrade.Task {
	logger.UpgradeLogger().Debugf("Marking task %s status %s", versionId, upgrade.InstallErr)
	doc, err := UpdateOrCreateTask(&upgrade.Task{
		VersionId:       versionId,
		Status:          upgrade.InstallErr,
		InstallStatus:   upgrade.Err,
		DoneInstallTime: time.Now().Format(time.RFC3339)})
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to mark task up err %s", err)
	}
	return doc
}

func MarkTaskInstalled(versionId string) (*upgrade.Task, error) {
	logger.UpgradeLogger().Debugf("Marking task %s status %s", versionId, upgrade.Installed)
	doc, err := UpdateOrCreateTask(&upgrade.Task{
		VersionId:     versionId,
		Status:        upgrade.Installed,
		InstallStatus: upgrade.Done,
		DoneDownTime:  time.Now().Format(time.RFC3339)})
	return doc, err
}
