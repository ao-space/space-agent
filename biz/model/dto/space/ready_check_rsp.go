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

package space

type ReadyCheckRsp struct {
	Paired             int  `json:"paired"`             // 0: 已经绑定; 1: 新盒子; 2: 已解绑
	DiskInitialCode    int  `json:"diskInitialCode"`    // 1: 磁盘正常; 2: 未初始化;3: 正在格式化; 4: 正在数据同步; 100: 未知错误; 101:  磁盘格式化错误; >101: 其他初始化错误;
	MissingMainStorage bool `json:"missingMainStorage"` // 缺少主存储
}
