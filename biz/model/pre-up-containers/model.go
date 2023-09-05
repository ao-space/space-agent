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

package preupcontainers

import (
	"agent/res"

	"agent/utils/logger"
	"github.com/dungeonsnd/gocom/encrypt/encoding"
)

type stPreUpContainers struct {
	Version         string   `json:"version"`
	PreUpContainers []string `json:"preUpContainers"`
}

var PreUpContainers *stPreUpContainers

func init() {
	content := res.GetContentPreUpContainers()
	if len(content) > 0 {
		r := &stPreUpContainers{}
		err := encoding.JsonDecode(content, r)
		if err != nil {
			logger.AppLogger().Warnf("failed GetContentPreUpContainers JsonDecode, err=%v", err)
		} else {
			PreUpContainers = r
		}
	}
}
