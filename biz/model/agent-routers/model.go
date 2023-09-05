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

package agentrouters

import (
	"agent/res"

	"agent/utils/logger"
	"github.com/dungeonsnd/gocom/encrypt/encoding"
)

type stAgentRouters struct {
	Version    string              `json:"version"`
	ValidPaths map[string][]string `json:"valid-paths"`
}

var AgentRouters *stAgentRouters

func init() {
	content := res.GetContentAgentRouters()
	if len(content) > 0 {
		r := &stAgentRouters{}
		err := encoding.JsonDecode(content, r)
		if err != nil {
			logger.AppLogger().Warnf("failed GetContentAgentRouters JsonDecode, err=%v", err)
		} else {
			AgentRouters = r
		}
	}
}
