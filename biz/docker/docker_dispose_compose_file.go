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
 * @Date: 2021-12-25 11:41:01
 * @LastEditors: jeffery
 * @LastEditTime: 2022-02-24 14:33:03
 * @Description:
 */

 package docker

 import (
	 "agent/biz/model/device"
	 "agent/biz/model/disk_initial/model"
	 "agent/config"
	 "agent/res"
	 "fmt"
	 math_rand "math/rand"
	 "strings"
	 "time"
 
	 "agent/utils/logger"
 
	 "github.com/dungeonsnd/gocom/file/fileutil"
 )
 
 func replacePlaceholderInComposeFile(f string) error {
	 logger.AppLogger().Debugf("replacePlaceholderInComposeFile, f:%v", f)
 
	 btid := device.GetDeviceInfo().Btid
 
	 b, err := fileutil.ReadFromFile(f)
	 if err != nil {
		 logger.AppLogger().Debugf("replacePlaceholderInComposeFile, failed ReadFromFile f:%v", f)
		 return err
	 }
	 s := string(b)
	 logger.AppLogger().Debugf("replacePlaceholderInComposeFile, %v", f)
	 s = strings.ReplaceAll(s, config.Config.Box.Loki.Placeholder, btid)
	 // logger.AppLogger().Debugf("replacePlaceholderInComposeFile, after Replace placeholder %v new content: %v", f, s)
	 newContent := []byte(s)
	 err = fileutil.WriteToFile(f, newContent, true)
	 if err != nil {
		 logger.AppLogger().Warnf("replacePlaceholderInComposeFile, failed WriteToFile %v", f)
		 return err
	 }
	 logger.AppLogger().Debugf("replacePlaceholderInComposeFile, succ WriteToFile %v", f)
	 return nil
 }
 
 func disposeComposeFile() {
	 f := config.Config.Docker.CustomComposeFile
	 logger.AppLogger().Debugf("disposeComposeFile, f:%v", f)
	 replacePlaceholderInComposeFile(f)
	 if fileutil.IsFileExist(f) {
		 logger.AppLogger().Infof("DisposeComposeFile, %v exist", f)
		 b, err := fileutil.ReadFromFile(f)
		 if err != nil {
			 logger.AppLogger().Warnf("failed ReadFromFile, file:%v, err:%v", f, err)
		 } else {
			 err1 := fileutil.WriteToFile(config.Config.Docker.ComposeFile,
				 b, true)
			 if err1 != nil {
				 logger.AppLogger().Warnf("failed WriteToFile, file:%v, err1:%v", config.Config.Docker.ComposeFile, err1)
			 }
		 }
	 } else {
		 logger.AppLogger().Warnf("DisposeComposeFile, %v not exist", f)
	 }
 
	 err := ProcessVolumes(config.Config.Docker.ComposeFile, model.GetFileStoragePath())
	 if err != nil {
		 logger.AppLogger().Warnf("failed ProcessVolumes, file:%v, err1:%v", config.Config.Docker.ComposeFile, err)
	 }
 }
 
 func writeDefaultDockerComposeFile() {
	 // rpm 安装时已经不释放了，所以须在程序里释放。升级时已经有该文件了，所以启动时根据配置项释放程序中内置的。
	 if fileutil.IsFileNotExist(config.Config.Docker.CustomComposeFile) {
		 fmt.Printf("writeDefaultDockerComposeFile, %v not exist\n", config.Config.Docker.CustomComposeFile)
 
		 composeFileContent := res.GetContentDockerCompose()
		 composeFileContent = replaceRandomPasswordAndPortPlaceholder(composeFileContent)
		 fileutil.WriteToFile(config.Config.Docker.CustomComposeFile, composeFileContent, true)
	 } else {
		 fmt.Printf("writeDefaultDockerComposeFile, %v exist\n", config.Config.Docker.CustomComposeFile)
 
		 if config.Config.OverwriteDockerCompose {
			 composeFileContent := res.GetContentDockerCompose()
			 composeFileContent = replaceRandomPasswordAndPortPlaceholder(composeFileContent)
			 fileutil.WriteToFile(config.Config.Docker.CustomComposeFile, composeFileContent, true)
		 }
	 }
 }
 
 func writeUpgradeComposeFile() {
	 f := config.Config.Docker.UpgradeComposeFile
	 logger.AppLogger().Debugf("dispose upgrade ComposeFile, f:%v", f)
	 composeFileContent := res.GetContentUpgradeComposeFile()
	 err := fileutil.WriteToFile(f, composeFileContent, true)
	 if err != nil {
		 logger.AppLogger().Errorf("write upgrade ComposeFile err:%v", err)
		 return
	 }
 }
 
 func replaceRandomPasswordAndPortPlaceholder(composeFileContent []byte) []byte {
	 logger.AppLogger().Debugf("replaceRandomPasswordAndPortPlaceholder")
 
	 randstr := ""
	 if fileutil.IsFileExist(config.Config.Box.RandDockercomposePassword) {
 
		 b, err := fileutil.ReadFromFile(config.Config.Box.RandDockercomposePassword)
		 if err != nil {
			 logger.AppLogger().Errorf("replaceRandomPasswordAndPortPlaceholder, ReadFromFile %v err: %v",
				 config.Config.Box.RandDockercomposePassword, err)
		 } else {
			 randstr = string(b)
		 }
	 }
	 if len(randstr) < 2 {
		 // 替换成随机密码
		 randstr = rand(16)
		 err := fileutil.WriteToFile(config.Config.Box.RandDockercomposePassword, []byte(randstr), true)
		 if err != nil {
			 logger.AppLogger().Errorf("replaceRandomPasswordAndPortPlaceholder, WriteToFile %v err: %v",
				 config.Config.Box.RandDockercomposePassword, err)
		 }
	 }
	 logger.AppLogger().Debugf("replaceRandomPasswordAndPortPlaceholder, randstr:%v", randstr)
	 composeFileContent = []byte(strings.ReplaceAll(string(composeFileContent), "placeholder_mysecretpassword", randstr))
 
	 randRedisPort := ""
	 if fileutil.IsFileExist(config.Config.Box.RandDockercomposeRedisPort) {
 
		 b, err := fileutil.ReadFromFile(config.Config.Box.RandDockercomposeRedisPort)
		 if err != nil {
			 logger.AppLogger().Errorf("replaceRandomPasswordAndPortPlaceholder, ReadFromFile %v err: %v",
				 config.Config.Box.RandDockercomposeRedisPort, err)
		 } else {
			 randRedisPort = string(b)
			 if err != nil {
				 logger.AppLogger().Errorf("replaceRandomPasswordAndPortPlaceholder, strconv.Atoi %v err: %v",
					 string(b), err)
			 }
		 }
	 }
	 if len(randRedisPort) < 2 {
		 // 替换成随机端口号
		 randRedisPort = randInt(19000, 19999) // redis 端口号在一个区间内随机生成
		 err := fileutil.WriteToFile(config.Config.Box.RandDockercomposeRedisPort, []byte(randRedisPort), true)
		 if err != nil {
			 logger.AppLogger().Errorf("replaceRandomPasswordAndPortPlaceholder, WriteToFile %v err: %v",
				 config.Config.Box.RandDockercomposeRedisPort, err)
		 }
	 }
 
	 i := strings.Index(string(composeFileContent), "placeholder_6379")
	 logger.AppLogger().Debugf("replaceRandomPasswordAndPortPlaceholder, randRedisPort:%v, strings.Index return :%v", randRedisPort, i)
 
	 composeFileContent = []byte(strings.ReplaceAll(string(composeFileContent), "placeholder_6379", randRedisPort))
 
	 // 更新自身连接 redis 的配置。
	 addrSplit := strings.Split(config.Config.Redis.Addr, ":")
	 if len(addrSplit) > 1 {
		 config.UpdateRedisConfig(addrSplit[0]+":6379", randstr)
	 } else {
		 config.UpdateRedisConfig("127.0.0.1:"+randRedisPort, randstr)
	 }
 
	 return composeFileContent
 }
 
 func rand(length int) string {
	 str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	 bytes := []byte(str)
	 result := []byte{}
	 r := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
	 for i := 0; i < length; i++ {
		 result = append(result, bytes[r.Intn(len(bytes))])
	 }
	 return string(result)
 }
 
 func randInt(min int, max int) string {
	 math_rand.Seed(time.Now().UnixNano())
	 r := math_rand.Intn(max-min) + min
	 return fmt.Sprintf("%v", r)
 }
 