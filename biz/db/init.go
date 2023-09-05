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
	"agent/config"
	"agent/utils/logger"
	scribble "github.com/nanobox-io/golang-scribble"
	"os"
	"path"
	"strings"
	"sync"
)

var lock sync.RWMutex

var conf = config.Config.RunTime
var Dir = path.Join(conf.BasePath, conf.DBDir)

func NewDBClient() (*scribble.Driver, error) {
	return scribble.New(Dir, nil)
}

func initDB(filePath string) error {
	defer lock.Unlock()
	lock.Lock()
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte("{}"))
	if err != nil {
		return err
	}
	f.Sync()
	defer f.Close()
	return nil
}

func CheckAndCreateDB() error {
	collPath := path.Join(Dir, conf.UpgradeCollection)
	err := os.MkdirAll(collPath, 0755)
	if err != nil {
		return err
	}
	taskPath := path.Join(collPath, conf.TaskResource+".json")
	_, err = os.Stat(taskPath)
	if err != nil {
		if !os.IsExist(err) {
			err := initDB(taskPath)
			if err != nil {
				return err
			}
		}
	}
	_, err = ReadTask("")
	if err != nil {
		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			logger.UpgradeLogger().Infof("As %v , will initialize and rewrite db file ", err)
			err = initDB(taskPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
