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

package did

import (
	"agent/biz/model/did/leveldb"
	"agent/utils/logger"
	"strings"
)

func SaveAoIdToDid(levelDBTrans *leveldb.Trans, aoId, did string) error {
	if i := strings.Index(did, "?"); i > 0 {
		did = did[:i]
	}
	if i := strings.Index(did, "#"); i > 0 {
		did = did[:i]
	}

	err := leveldb.Put(levelDBTrans, []byte(leveldb.KNameOfAoIdToDid(aoId)), []byte(did))
	if err != nil {
		logger.AppLogger().Debugf("SaveAoIdToDid, leveldb.Put:%v, err:%v",
			leveldb.KNameOfAoIdToDid(aoId), err)
		return err
	}

	err = leveldb.Put(levelDBTrans, []byte(leveldb.KNameOfDidToAoId(did)), []byte(aoId))
	if err != nil {
		logger.AppLogger().Debugf("SaveAoIdToDid, leveldb.Put:%v, err:%v",
			leveldb.KNameOfAoIdToDid(did), err)
		return err
	}

	return nil
}

func GetDidByAoId(levelDBTrans *leveldb.Trans, aoId string) (string, bool, error) {

	exist, err := leveldb.Has(levelDBTrans, []byte(leveldb.KNameOfAoIdToDid(aoId)))
	if err != nil { // error
		logger.AppLogger().Debugf("GetDidByAoId, leveldb.Has:%v, err:%v",
			leveldb.KNameOfAoIdToDid(aoId), err)
		return "", false, err
	}

	if exist {
		didStr, err := leveldb.Get(levelDBTrans, []byte(leveldb.KNameOfAoIdToDid(aoId)))
		if err != nil { // error
			logger.AppLogger().Debugf("GetDidByAoId, leveldb.Get:%v, err:%v",
				leveldb.KNameOfAoIdToDid(aoId), err)
			return "", false, err
		}

		return string(didStr), true, err
	} else {
		return "", false, nil
	}
}

func GetAoIdByDid(levelDBTrans *leveldb.Trans, did string) (string, bool, error) {
	if i := strings.Index(did, "?"); i > 0 {
		did = did[:i]
	}
	if i := strings.Index(did, "#"); i > 0 {
		did = did[:i]
	}

	exist, err := leveldb.Has(levelDBTrans, []byte(leveldb.KNameOfDidToAoId(did)))
	if err != nil { // error
		logger.AppLogger().Debugf("GetAoIdByDid, leveldb.Has:%v, err:%v",
			leveldb.KNameOfDidToAoId(did), err)
		return "", false, err
	}

	if exist {
		aoId, err := leveldb.Get(levelDBTrans, []byte(leveldb.KNameOfDidToAoId(did)))
		if err != nil { // error
			logger.AppLogger().Debugf("GetAoIdByDid, leveldb.Get:%v, err:%v",
				leveldb.KNameOfDidToAoId(did), err)
			return "", false, err
		}

		return string(aoId), false, err
	} else {
		return "", false, nil
	}
}
