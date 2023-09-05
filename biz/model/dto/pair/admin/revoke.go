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
 * @Date: 2021-11-24 15:33:06
 * @LastEditors: wenchao
 * @LastEditTime: 2021-12-22 14:31:58
 * @Description:
 */

package admin

type RevokReq struct {
	Password string `json:"password"`
}

// type RevokResult struct {
// 	ErrorTimes      int    `json:"errorTimes"`      // 已经尝试错误次数
// 	LeftTryTimes    int    `json:"leftTryTimes"`    // 剩余次数次数
// 	TryAfterSeconds int    `json:"tryAfterSeconds"` // 多少秒后重试
// 	BoxUuid         string `json:"boxUuid"`         // 盒子端唯一id.
// }
