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
 * @LastEditors: jeffery
 * @LastEditTime: 2022-06-06 11:28:46
 * @Description:
 */
package main

import (
	"agent/biz/alivechecker"
	"agent/biz/disk_space_monitor/log_dir_monitor"
	"agent/biz/docker"
	"agent/config"
	"strings"

	"agent/biz/model/clientinfo"
	"agent/biz/model/device"
	"agent/biz/model/did/leveldb"
	"agent/biz/service/platform"
	"agent/biz/service/upgrade"
	"agent/biz/web"
	"agent/utils/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	serviceswithplatform "agent/biz/service/switch-platform"
	// _ "net/http/pprof"
)

func main() {

	logger.SetLogPath(config.Config.Log.Path)
	logger.SetLogConfig(int(config.Config.Log.RotationSize),
		int(config.Config.Log.RotationCount),
		int(config.Config.Log.MaxAge), false)
	logger.PrecreateAllLoggers()

	config.Version = Version
	config.VersionNumber = VersionNumber

	logger.AppLogger().Infof("================[%v Started] [system-agent version:%v]================",
		os.Args[0], config.Version+"-"+config.VersionNumber)

	fmt.Printf("AoLogDirBase:%+v\n", config.Config.Log.AoLogDirBase)
	fmt.Printf("singleDockerModeEnv:%+v\n", os.Getenv(config.Config.Box.RunInDocker.AoSpaceSingleDockerModeEnv))

	if err := AgentCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	device.InitDeviceInfo()
	device.InitDeviceKey()
	clientinfo.InitClientInfo()

	if !strings.EqualFold(os.Getenv(config.Config.Box.RunInDocker.AoSpaceSingleDockerModeEnv), "true") {
		go platform.InitPlatformAbility()
		serviceswithplatform.RetryUnfinishedStatus()
		upgrade.CronForUpgrade()
	}

	// 启动 web/http api 服务
	web.Start()

	// 启动 docker 微服务创建或启动
	if strings.EqualFold(os.Getenv(config.Config.Box.RunInDocker.AoSpaceSingleDockerModeEnv), "true") {
		docker.MigrateFileStorageData()
	} else {
		docker.Start()
	}
	alivechecker.Start()

	// 检测是否需要发送升级推送
	if !strings.EqualFold(os.Getenv(config.Config.Box.RunInDocker.AoSpaceSingleDockerModeEnv), "true") {
		go upgrade.CheckUpgradeSucc()
	}

	// 日志目录监控
	log_dir_monitor.Start()

	if err := leveldb.OpenDB(); err != nil {
		fmt.Printf("\nFailed leveldb.OpenDB, err:%v\n", err)
		os.Exit(0)
	}

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	for s := range quitChan {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Printf("signal %v\n", s)
			GracefullExit()
		case syscall.SIGUSR1:
			// fmt.Println("usr1 signal", s)
		case syscall.SIGUSR2:
			// fmt.Println("usr2 signal", s)
		default:
			// fmt.Println("other signal", s)
		}
	}

}

func GracefullExit() {
	os.Exit(0)
}
