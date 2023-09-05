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
 * @Date: 2021-12-02 16:55:11
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-02 22:50:11
 * @Description:
 */

package storage

import (
	"sync"

	"github.com/dungeonsnd/gocom/file/fileutil"
)

type Storage struct {
	Lock     sync.Mutex
	FileName string
}

func NewStorage(fileName string) *Storage {
	return &Storage{Lock: sync.Mutex{}, FileName: fileName}
}

func (store *Storage) SaveJson(obj interface{}) error {
	store.Lock.Lock()
	defer store.Lock.Unlock()
	return fileutil.WriteToFileAsJson(store.FileName, obj, "  ", true)
}

func (store *Storage) LoadJson(obj interface{}) error {
	store.Lock.Lock()
	defer store.Lock.Unlock()
	return fileutil.ReadFileJsonToObject(store.FileName, obj)
}
