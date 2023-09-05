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
	"agent/utils/logger"
)

func Put(levelDBTrans *Trans, key, value []byte) error {
	logger.LevelDBLogger().Debugf("[Before] leveldb.Put, key:%v, value:%v",
		string(key), string(value))
	if levelDBTrans != nil {
		return levelDBTrans.Put(key, value, nil)
	}
	logger.LevelDBLogger().Debugf("leveldb.Put, key:%v, value:%v",
		string(key), string(value))
	return ldb.Put(key, value, nil)
}

func Get(levelDBTrans *Trans, key []byte) ([]byte, error) {
	logger.LevelDBLogger().Debugf("[Before] leveldb.Get, key:%v",
		string(key))
	if levelDBTrans != nil {
		return levelDBTrans.Get([]byte(key), nil)
	}
	v, err := ldb.Get([]byte(key), nil)
	logger.LevelDBLogger().Debugf("leveldb.Get, key:%v, value:%v, err:%v",
		string(key), string(v), err)
	return v, err
}

func Has(levelDBTrans *Trans, key []byte) (bool, error) {
	logger.LevelDBLogger().Debugf("[Before] leveldb.Has, key:%v",
		string(key))
	if levelDBTrans != nil {
		return levelDBTrans.Has([]byte(key), nil)
	}
	exist, err := ldb.Has([]byte(key), nil)
	logger.LevelDBLogger().Debugf("leveldb.Has, key:%v, exist:%v, err:%v",
		string(key), exist, err)
	return exist, err
}

func Delete(levelDBTrans *Trans, key []byte) error {
	logger.LevelDBLogger().Debugf("[Before] leveldb.Delete, key:%v",
		string(key))
	if levelDBTrans != nil {
		return levelDBTrans.Delete([]byte(key), nil)
	}
	err := ldb.Delete([]byte(key), nil)
	logger.LevelDBLogger().Debugf("leveldb.Delete, key:%v, err:%v",
		string(key), err)
	return err
}
