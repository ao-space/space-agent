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

package upgrade

import (
	"agent/biz/db"
	"agent/biz/model/dto"
	upModel "agent/biz/model/upgrade"
	"agent/biz/service/upgrade"
	"agent/config"
	"agent/utils/hardware"
	"fmt"
	"github.com/dungeonsnd/gocom/log4go"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"agent/utils/logger"
	"github.com/gin-gonic/gin"
)

var conf = config.Config.RunTime

// 下载和安装都是异步的， 但是如果中间失败，不重试。 直接标记状态。

// GetUpgradeConfig godoc
// @Summary get current upgrade config
// @Tags upgrade
// @Produce  application/json
// @Success 200 {object} upModel.UpgradeConfig
// @Router /agent/v1/api/upgrade/config [GET]
func GetUpgradeConfig(c *gin.Context) {
	// uConf := config.Config.Box.UpgradeConfig
	var upConf *upModel.UpgradeSettings
	if _, err := os.Stat(config.Config.Box.UpgradeConfig.SettingsFile); err != nil {
		settings := &upModel.UpgradeSettings{AutoDownload: true, AutoInstall: true}
		upModel.SetUpgradeSettings(settings)
		if err != nil {
			logger.AppLogger().Warnf("failed SetUpgradeSettings, file:%v, err:%v",
				config.Config.Box.UpgradeConfig.SettingsFile, err)
			return
		}
		logger.AppLogger().Debugf("succ SetUpgradeSettings, file:%v, settings:%+v",
			config.Config.Box.UpgradeConfig.SettingsFile, settings)
	}
	upConf = upModel.GetUpgradeSettings()
	resp := upModel.UpgradeConfig{AutoDownload: upConf.AutoDownload,
		AutoInstall: upConf.AutoInstall}
	c.JSON(http.StatusOK, resp)
}

// SetUpgradeConfig godoc
// @Summary set auto upgrade config
// @Tags upgrade
// @Param  config body upModel.UpgradeConfig true "upgrade config"
// @Accept   json
// @Produce   json
// @Failure 400 string dto.BaseRsp
// @Success 200 {object} upModel.UpgradeConfig
// @Router /agent/v1/api/upgrade/config [POST]
func SetUpgradeConfig(c *gin.Context) {

	logger.AppLogger().Infof("/agent/v1/api/upgrade/config [system-agent version:%v]", config.Version)

	upConf := upModel.UpgradeConfig{}
	err := c.BindJSON(&upConf)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseRsp{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	err = upModel.SetUpgradeSettings(&upModel.UpgradeSettings{AutoDownload: upConf.AutoDownload,
		AutoInstall: upConf.AutoInstall})
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BaseRsp{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, config.Config.Box.UpgradeConfig)
}

// GetTaskStatus godoc
// @Summary get upgrade task status
// @Tags upgrade
// @Produce   json
// @Success 200 {object} upModel.Task
// @Failure 400 string dto.BaseRsp  "there is no task"
// @Failure 500 string dto.BaseRsp  "please try after seconds"
// @Router /agent/v1/api/upgrade/status [GET]
func GetTaskStatus(c *gin.Context) {
	task, err := db.ReadTask("")
	if err != nil {
		if task.VersionId == "" {
			c.JSON(http.StatusBadRequest, dto.BaseRsp{Code: http.StatusBadRequest, Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.BaseRsp{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	if task.Status == upModel.Downloaded {
		lastVersion, err := upgrade.CheckLatestVersion()
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.BaseRsp{Code: http.StatusInternalServerError, Message: err.Error()})
			return
		}
		if lastVersion.PkgVersion != task.VersionId {
			upgrade.RecheckUpgradeStatus(false)
		}
	}
	defer logger.AccessLogger().Infof("get upgrade status,%v", task)
	c.JSON(http.StatusOK, task)
}

// StartDownload godoc
// @Summary add a new download task
// @Tags upgrade
// @Accept json
// @Produce   json
// @Param down body upModel.StartDownRes true "download params"
// @Success 200 {object} upModel.Task
// @Success 208 {object} upModel.Task  "new server version is downloaded,but not installed"
// @Failure 400 string dto.BaseRsp  "new version does not exist"
// @Failure 500 string dto.BaseRsp  "failed , please try again"
// @Router /agent/v1/api/upgrade/download [POST]
func StartDownload(c *gin.Context) {
	logger.UpgradeLogger().Infof("start: /agent/v1/api/upgrade/download")

	resource := "task"
	downArgs := upModel.StartDownRes{}
	err := c.ShouldBindJSON(&downArgs)
	if err != nil {

		logger.UpgradeLogger().Errorf("Failed to bind args %s", err)
		c.JSON(http.StatusBadRequest, dto.BaseRsp{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	overallInfo, err := upgrade.GetLatestVersionMetadata()
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to download latest pkg docker-compose: %s", err)
		c.JSON(http.StatusInternalServerError, dto.BaseRsp{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	logger.UpgradeLogger().Debugf("overallInfo,%v", overallInfo)
	if overallInfo.VersionId != downArgs.VersionId {
		err := fmt.Errorf("%s It's not lastest version, the lastest version is %s", downArgs.VersionId, overallInfo.VersionId)
		//logger.AccessLogger()
		c.JSON(http.StatusBadRequest, dto.BaseRsp{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	//versionDesc, err := upgrade.CheckLatestVersion()

	dbClient, err := db.NewDBClient()
	task := new(upModel.Task)
	err = dbClient.Read(conf.UpgradeCollection, conf.TaskResource, &task)
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to Read db %v", err)
		c.JSON(http.StatusInternalServerError, dto.BaseRsp{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	if task.VersionId == overallInfo.VersionId && task.Status == upModel.Downloaded && !downArgs.Anew {
		mess := fmt.Sprintf("This version %s is already downloaded", overallInfo.VersionId)
		logger.UpgradeLogger().Warnf("Ignored: " + mess)
		c.JSON(
			http.StatusAlreadyReported,
			dto.BaseRsp{Code: http.StatusAlreadyReported, Message: mess})
		return
	}
	if task.VersionId == overallInfo.VersionId && task.Status == upModel.Downloading && !downArgs.Anew {
		mess := fmt.Sprintf("This version %s is %s, please wait", overallInfo.VersionId, upModel.Downloading)
		logger.UpgradeLogger().Warnf("Ignored: " + mess)
		c.JSON(
			http.StatusAlreadyReported,
			dto.BaseRsp{Code: http.StatusAlreadyReported, Message: mess})
		return
	}

	task.VersionId = overallInfo.VersionId
	task.Status = upModel.Downloading
	task.DownStatus = upModel.Ing
	task.StartDownTime = time.Now().Format(time.RFC3339)
	task.NeedReboot = overallInfo.Restart
	if overallInfo.KernelUrl != "" && overallInfo.KernelVersion != "" {
		task.KernelImg.VersionId = overallInfo.KernelVersion
		task.KernelImg.PkgPath = upgrade.OTAImagePath
	}
	err = dbClient.Write(conf.UpgradeCollection, resource, task)
	if err != nil {
		logger.UpgradeLogger().Errorf("write file db error:%v", err)
		c.JSON(
			http.StatusInternalServerError,
			dto.BaseRsp{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	var agent upgrade.Agent
	if hardware.RunningInDocker() {
		agent = &upgrade.ContainerAgent{
			VersionId:   task.VersionId,
			ComposeFile: filepath.Join(config.Config.RunTime.BasePath, config.Config.RunTime.PkgDir, "docker-compose.yml"),
		}
	} else {
		agent = &upgrade.NativeAgent{VersionId: task.VersionId}
	}
	go agent.Download()
	c.JSON(http.StatusOK, task)
}

// StartUpgrade godoc
// @Summary POST:/upgrade/start-upgrade start to install new version
// @Tags upgrade
// @Accept   json
// @Produce   json
// @Param upgrade body upModel.StartUpgradeRes true "install params"
// @Success 200 {object} upModel.Task
// @Failure 400 string dto.BaseRsp  "version have not been downloaded"
// @Failure 409 string dto.BaseRsp  "there is task which is running"
// @Failure 500 string dto.BaseRsp  "failed ,please try again"
// @Router /agent/v1/api/upgrade/install [POST]
func StartUpgrade(c *gin.Context) {
	log4go.I("/agent/v1/api/upgrade/install")

	resData := upModel.StartUpgradeRes{}
	err := c.ShouldBindJSON(&resData)
	if err != nil {
		logger.UpgradeLogger().Errorf("Failed to bind args %s", err)
		c.JSON(http.StatusBadRequest, dto.BaseRsp{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}
	versionId := strings.TrimSpace(resData.VersionId)
	errCode, err := upgrade.CheckVersionStatusRight(versionId)
	if err != nil {
		c.JSON(errCode, dto.BaseRsp{Code: errCode, Message: err.Error()})
		return
	}
	task, err := db.MarkTaskInstalling(versionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BaseRsp{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}
	var agent upgrade.Agent
	c.JSON(http.StatusOK, task)
	go func() {
		// 升级agent
		// 升级agent
		if hardware.RunningInDocker() {
			agent = &upgrade.ContainerAgent{
				VersionId: versionId}
			go agent.Upgrade()
		} else {
			upgrade.InstallAgentV2(versionId)
		}
	}()

}
