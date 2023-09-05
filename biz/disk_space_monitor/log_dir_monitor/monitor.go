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

package log_dir_monitor

import (
	"agent/config"
	"os"
	"time"

	"agent/utils/logger"

	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/robfig/cron"
)

func Start() {
	cronTable()
}

func cronTable() {
	c := cron.New()
	err := c.AddFunc(config.Config.Log.AoLogDirCheckCronExp, checkSize)
	if err != nil {
		logger.AppLogger().Warnf("Failed to config cron: %s", err)
		return
	}
	c.Start()
}

func checkSize() {
	d := config.Config.Log.AoLogDirBase
	if fileutil.IsFileNotExist(d) {
		return
	}
	sz, err := fileutil.DirSize(d)
	if err != nil {
		logger.AppLogger().Warnf("Failed to checkSize, err: %v", err)
		return
	}

	if sz < config.Config.Log.AoLogDirBaseSizLimit {
		logger.AppLogger().Debugf("checkSize, return, sz: %v < BpLogDirBaseSizLimit:%v", sz, config.Config.Log.AoLogDirBaseSizLimit)
		return
	}
	logger.AppLogger().Debugf("checkSize, sz: %v >= BpLogDirBaseSizLimit:%v", sz, config.Config.Log.AoLogDirBaseSizLimit)

	files, err := fileutil.ListDir(d)
	if err != nil {
		logger.AppLogger().Warnf("Failed to ListDir, err: %v", err)
		return
	}
	logger.AppLogger().Debugf("len(files): %v", len(files))

	for _, logfile := range files {
		st, err := os.Stat(logfile)
		if err != nil {
			logger.AppLogger().Warnf("Failed to Stat, err: %v", err)
			continue
		}

		if time.Now().Local().After(st.ModTime().Add(time.Hour * 24 * time.Duration(config.Config.Log.AoLogMaxDayLimit))) {
			err1 := os.Remove(logfile)
			if err1 != nil {
				logger.AppLogger().Warnf("Failed to Remove, err: %v", err1)
				continue
			} else {
				logger.AppLogger().Debugf("Remove logfile: %v", logfile)
			}
		}
	}

}
