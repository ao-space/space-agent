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

package logger

import (
	zaplogger "agent/deps/logger"

	"go.uber.org/zap"
)

var (
	path                = ""
	LogFileApp          = "app.log"
	LogFileCheck        = "check.log"
	LogFileAccess       = "access.log"
	LogFileNotification = "notification.log"
	LogFileLedStatus    = "ledstatus.log"
	LogFileIpScan       = "ipscan.log"
	LogFileUpgrade      = "upgrade.log"
	LogFileCertificate  = "certificate.log"
	LogFileDocker       = "docker.log"
	LogFileLevelDB      = "leveldb.log"
)

func SetLogPath(p string) {
	path = p
}

func SetLogConfig(MaxSize, MaxBackups, MaxAge int, Compress bool) {
	zaplogger.SetLogConfig(&zaplogger.LogConfig{MaxSize: MaxSize,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
		Compress:   Compress})
}

// "debug", "info", "warn", "error", "dpanic", "panic", and "fatal"
func SetLevel(level string) {
	zaplogger.SetLevel(level)
}

// 注意! 所有日志都需要提前初始化。
//
//	日志模块为了性能, 读写 map 没加锁。
//	所以需要提前初始化，初始化以后只会读取 map ，不会写入了。并发读取是没问题的。
func PrecreateAllLoggers() {
	AppLogger()
	CheckLogger()
	AccessLogger()
	NotificationLogger()
	LedStatusLogger()
	IpScanLogger()
	CertificateLogger()
	DockerLogger()
	LevelDBLogger()
}

func AppLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileApp)
}

func UpgradeLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileUpgrade)
}

func CheckLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileCheck)
}

func AccessLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileAccess)
}

func NotificationLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileNotification)
}

func LedStatusLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileLedStatus)
}

func IpScanLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileIpScan)
}

func CertificateLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileCertificate)
}

func DockerLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileDocker)
}

func LevelDBLogger() *zap.SugaredLogger {
	return zaplogger.Logger(path + LogFileLevelDB)
}
