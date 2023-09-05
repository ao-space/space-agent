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

package routers

import (
	"agent/biz/docker"
	"agent/biz/model/device_ability"
	"agent/config"
	"agent/utils/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"os"
)

type ExternalWebServer struct {
	Router *gin.Engine
}

type InternalWebServer struct {
	Router *gin.Engine
}

type HttpServer interface {
	Start()
}

// Start external web server start
func (w *ExternalWebServer) Start() {
	fmt.Printf("startWebServer \n")
	logger.AppLogger().Infof("startWebServer")

	if gin.Mode() == gin.DebugMode {
		w.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	err := w.Router.Run(config.Config.Web.DefaultListenAddr)
	if err != nil {
		logger.AppLogger().Errorf("Failed startWebServer, err: %v", err)
		os.Exit(-1)
	}
	return
}

// Start internal web server start
func (w *InternalWebServer) Start() {
	fmt.Printf("startWebServerDockerLocal \n")
	logger.AppLogger().Infof("startWebServerDockerLocal")

	if gin.Mode() == gin.DebugMode {
		w.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	localListenAddr := config.Config.Web.DockerLocalListenAddr
	abilityModel := device_ability.GetAbilityModel()
	if abilityModel.RunInDocker {
		localListenAddr = config.Config.Web.DockerLocalListenAddrRunInDocker
	}

	logger.AppLogger().Debugf("startWebServerDockerLocal, using %v ", localListenAddr)
	err := w.Router.Run(localListenAddr)
	if err != nil {
		err1 := fmt.Errorf("Failed startWebServerDockerLocal using %v, err: %v", localListenAddr, err)
		fmt.Printf("%+v\n", err1)
		logger.AppLogger().Errorf("%+v", err1)
		// os.Exit(0) // 可能不是致命的
	} else {
		docker.UnsubscribeDockerNetwork(nil)
		return
	}
}
