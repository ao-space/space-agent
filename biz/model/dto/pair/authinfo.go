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
 * @Date: 2021-11-10 14:41:01
 * @LastEditors: wenchao
 * @LastEditTime: 2021-11-25 10:35:46
 * @Description:
 */
package pair

type AuthInfo struct {
	AuthKey string `json:"authKey"` // 盒子端授权给客户端的凭证. 使用客户端的公钥加密过了.
}
