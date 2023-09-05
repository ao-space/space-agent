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

package docker

import (
	preupcontainers "agent/biz/model/pre-up-containers"
	"agent/config"
	"agent/utils/tools"

	"agent/utils/logger"
)

func dockerPull() {
	setDockerStatus(ContainersDownloading)
	logger.AppLogger().Debugf("dockerPull Begin")
	err := docker.Pull(config.Config.Docker.ComposeFile)
	if err != nil {
		logger.AppLogger().Warnf("failed docker.Pull, err:%v ", err)
		setDockerStatus(ContainersDownloadedFail)
		return
	} else {
		logger.AppLogger().Debugf("SUCC docker.Pull")
	}
	setDockerStatus(ContainersDownloaded)
	logger.AppLogger().Debugf("@@ dockerPull Finished")
}

func dockerCreate() {
	logger.AppLogger().Debugf("dockerCreate Begin")
	setDockerStatus(ContainersStarting)
	err := docker.Create(config.Config.Docker.ComposeFile)
	if err != nil {
		logger.AppLogger().Warnf("failed docker.Create, err:%v ", err)
		return
	} else {
		logger.AppLogger().Debugf("SUCC docker.Create")
	}
	logger.AppLogger().Debugf("@@ dockerCreate Finished")
}

func dockerStart() {
	logger.AppLogger().Debugf("dockerStart Begin")
	setDockerStatus(ContainersStarting)
	err := docker.Start(config.Config.Docker.ComposeFile)
	if err != nil {
		logger.AppLogger().Warnf("failed docker.Start, err:%v ", err)
		setDockerStatus(ContainersStartedFail)
		return
	} else {
		logger.AppLogger().Debugf("SUCC docker.Start")
		setDockerStatus(ContainersStarted)
	}
	logger.AppLogger().Debugf("@@ dockerStart Finished")
}

func dockerDown() error {
	setDockerStatus(ContainersUnStarted)

	logger.AppLogger().Debugf("dockerDown Begin")
	err := docker.DownContainers(config.Config.Docker.ComposeFile)
	if err != nil {
		logger.AppLogger().Warnf("failed docker.DownContainers, err:%v ", err)
		setDockerStatus(ContainersDownloadedFail)
		return err
	} else {
		logger.AppLogger().Debugf("SUCC docker.DownContainers")
	}
	logger.AppLogger().Infof("@@ dockerDown Finished")
	return nil
}

func dockerStop(containerName string) error {
	logger.AppLogger().Debugf("dockerStop Begin, containerName:%v", containerName)
	err := docker.StopSpecifiedContainers(config.Config.Docker.ComposeFile, containerName)
	if err != nil {
		logger.AppLogger().Warnf("failed docker.StopSpecifiedContainers, err:%v ", err)
		return err
	} else {
		logger.AppLogger().Debugf("SUCC docker.StopSpecifiedContainers")
	}
	logger.AppLogger().Infof("@@ dockerStop Finished")
	return nil
}

func ContainersUpAndPrune(composeFile string, excludeServices []string) error {
	setDockerStatus(ContainersStarting)
	logger.AppLogger().Debugf("dockerUpAndPrune Begin")
	_, stdErr, err := docker.UpContainers(composeFile, excludeServices)
	if err != nil {
		logger.AppLogger().Warnf("Failed docker.UpContainers, err:%v ", err)
		setDockerStatus(ContainersStartedFail)
		return err
	} else if len(stdErr) != 0 {
		logger.AppLogger().Warnf("Failed docker UpContainers, err:%v ", stdErr)
		setDockerStatus(ContainersStartedFail)
		return err
	} else {
		logger.AppLogger().Debugf("@@ SUCC docker  UpContainers")
		setDockerStatus(ContainersStarted)
		go RemoveOldImage()
	}
	logger.AppLogger().Infof("@@ dockerUpAndPrune Finished")
	return nil
}

// func dockerUpAndPrune() {
// 	setDockerStatus(ContainersStarting)
// 	logger.AppLogger().Debugf("dockerUpAndPrune Begin")
// 	_, stdErr, err := docker.UpContainers(config.Config.Docker.ComposeFile)
// 	if err != nil {
// 		logger.AppLogger().Warnf("Failed docker.UpContainers, err:%v ", err)
// 		setDockerStatus(ContainersStartedFail)
// 		return
// 	} else if len(stdErr) != 0 {
// 		logger.AppLogger().Warnf("Failed docker UpContainers, err:%v ", stdErr)
// 		setDockerStatus(ContainersStartedFail)
// 		return
// 	} else {
// 		logger.AppLogger().Debugf("@@ SUCC docker  UpContainers")
// 		setDockerStatus(ContainersStarted)
// 		go RemoveOldImage()
// 	}
// 	logger.AppLogger().Infof("@@ dockerUpAndPrune Finished")
// }

func dockerUpAndPruneWithNoRecreate(excludeServices []string) {
	setDockerStatus(ContainersStarting)
	startProgress()
	logger.AppLogger().Debugf("dockerUpAndPruneWithNoRecreate Begin")
	_, stdErr, err := docker.UpContainersWithNoRecreate(config.Config.Docker.ComposeFile, excludeServices)
	if err != nil {
		logger.AppLogger().Warnf("Failed docker.UpContainersWithNoRecreate, err:%v ", err)
		setDockerStatus(ContainersStartedFail)
		stopProgress()
		return
	} else if len(stdErr) != 0 {
		logger.AppLogger().Warnf("Failed docker UpContainersWithNoRecreate, err:%v ", stdErr)
		setDockerStatus(ContainersStartedFail)
		stopProgress()
		return
	} else {
		logger.AppLogger().Debugf("@@ SUCC docker  UpContainersWithNoRecreate")
		// setDockerStatus(ContainersStarted) // 改成了检测网关接口了.
		go RemoveOldImage()
	}
	finishProgress()
	logger.AppLogger().Infof("@@ dockerUpAndPruneWithNoRecreate Finished")
}

func dockerPreUp() {
	setDockerPreUp(ContainersStarting)
	logger.AppLogger().Debugf("dockerPreUp Begin")
	// docker-compose -f docker-compose.yml up -d \
	// monitor-prometheus monitor-nodeexporter monitor-dockerexporter monitor-promtail \
	// aospace-postgresql aospace-rabbitmq aospace-filepreview aospace-nginx

	_, stdErr, err := docker.UpSpecifiedContainers(config.Config.Docker.ComposeFile,
		preupcontainers.PreUpContainers.PreUpContainers...)
	if err != nil {
		logger.AppLogger().Warnf("dockerPreUp, Failed docker.UpSpecifiedContainers, err:%v ", err)
		setDockerPreUp(ContainersStartedFail)
		return
	} else if len(stdErr) != 0 {
		logger.AppLogger().Warnf("dockerPreUp, Failed docker UpSpecifiedContainers, err:%v ", stdErr)
		setDockerPreUp(ContainersStartedFail)
		return
	} else {
		logger.AppLogger().Debugf("@@ dockerPreUp, SUCC docker  UpSpecifiedContainers")
		setDockerPreUp(ContainersStarted)
	}
	logger.AppLogger().Infof("@@ dockerPreUp Finished")
}

// RemoveOldImage is used to remove old versions of images
func RemoveOldImage() {
	logger.AppLogger().Debugf("Removing old versions of docker images")
	stdOut, stdErr, err := tools.ExeCmd("docker", "system", "prune", "-a", "-f")
	if err != nil || stdErr != "" {
		logger.AppLogger().Warnf("Failed to remove old versions docker image, stdOut:%v, stdErr:%v, err:%v",
			stdOut, stdErr, err)
	}
	logger.AppLogger().Debugf("Removing old versions of docker images return")
}
