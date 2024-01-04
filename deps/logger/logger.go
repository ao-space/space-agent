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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxsize" yaml:"maxsize"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxage" yaml:"maxage"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress" yaml:"compress"`
}

var logConfig *LogConfig
var atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

func SetLogConfig(config *LogConfig) {
	logConfig = config
}

var mLogger map[string]*zap.SugaredLogger

var mDefaultLoggerPath string

func init() {
	mLogger = make(map[string]*zap.SugaredLogger)
	mDefaultLoggerPath = "./"
}

func SetDefaultLoggerPath(defaultLoggerPath string) {
	mDefaultLoggerPath = defaultLoggerPath
}

func DefaultLogger() *zap.SugaredLogger {
	return Logger(mDefaultLoggerPath + "app.log")
}

func Logger(file string) *zap.SugaredLogger {
	log, ok := mLogger[file]
	if !ok {
		log = initLogger(file, atomicLevel)
		mLogger[file] = log
	}
	return log
}

// "debug", "info", "warn", "error", "dpanic", "panic", and "fatal"
func SetLevel(level string) {
	atomicLevel.UnmarshalText([]byte(level))
}

func initLogger(file string, level zapcore.LevelEnabler) *zap.SugaredLogger {
	encoder := getEncoder()
	writeSyncer := getLogWriter(file)
	core := zapcore.NewCore(encoder, writeSyncer, level)

	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n\n",
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(file string) zapcore.WriteSyncer {
	// file, _ := os.Create("./test.log")
	// return zapcore.AddSync(file)

	if logConfig != nil {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   file,
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		}
		return zapcore.AddSync(lumberJackLogger)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    20,
		MaxBackups: 4,
		MaxAge:     90,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
