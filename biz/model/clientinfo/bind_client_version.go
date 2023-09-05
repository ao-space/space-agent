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

package clientinfo

import "sync"

var mapClientVersion map[string]string
var lock sync.Mutex

func init() {
	mapClientVersion = make(map[string]string)
}

func SetClientVersion(k, v string) {
	lock.Lock()
	defer lock.Unlock()
	mapClientVersion[k] = v
}

func GetClientVersion(k string) (string, bool) {
	lock.Lock()
	defer lock.Unlock()
	v, found := mapClientVersion[k]
	return v, found
}
