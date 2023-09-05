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

package leveldb

import (
	"agent/config"
	"agent/utils/logger"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

var ldb *leveldb.DB

func OpenDB() error {
	f := fmt.Sprintf("%v/%v", config.Config.Box.DID.RootPath,
		config.Config.Box.DID.DBFileName)
	var err error
	ldb, err = leveldb.OpenFile(f, nil)
	if err != nil {
		logger.AppLogger().Warnf("OpenDB, file:%+v, err:%+v", f, err)
	}
	return err
}

func CloseDB() {
	logger.LevelDBLogger().Debugf("leveldb.Close, ldb:%+v", ldb)
	if ldb != nil {
		ldb.Close()
		ldb = nil
	}
}
