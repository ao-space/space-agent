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

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Trans struct {
	Transaction *leveldb.Transaction
}

func BeginTransaction() (*Trans, error) {
	logger.LevelDBLogger().Debugf("BeginTransaction")
	trans := &Trans{}

	var err error
	trans.Transaction, err = ldb.OpenTransaction()
	if err != nil {
		logger.LevelDBLogger().Warnf("BeginTrans, OpenTransaction err:%+v", err)
		return nil, err
	}
	return trans, nil
}

func (t *Trans) Commit() error {
	return t.Transaction.Commit()
}

func (t *Trans) Rollback() {
	t.Transaction.Discard()
}

func (t *Trans) Put(key, value []byte, wo *opt.WriteOptions) error {
	logger.LevelDBLogger().Debugf("leveldb.Put, key:%v, value:%v",
		string(key), string(value))
	return t.Transaction.Put(key, value, wo)
}

func (t *Trans) Get(key []byte, wo *opt.ReadOptions) ([]byte, error) {
	v, err := t.Transaction.Get([]byte(key), wo)
	logger.LevelDBLogger().Debugf("leveldb.Get, key:%v, value:%v, err:%v",
		string(key), string(v), err)
	return v, err
}

func (t *Trans) Has(key []byte, wo *opt.ReadOptions) (bool, error) {
	exist, err := t.Transaction.Has([]byte(key), wo)
	logger.LevelDBLogger().Debugf("leveldb.Has, key:%v, exist:%v, err:%v",
		string(key), exist, err)
	return exist, err
}

func (t *Trans) Delete(key []byte, wo *opt.WriteOptions) error {
	err := t.Transaction.Delete([]byte(key), wo)
	logger.LevelDBLogger().Debugf("leveldb.Delete, key:%v, err:%v",
		string(key), err)
	return err
}
