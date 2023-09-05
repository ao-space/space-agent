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
 * @Author: zhongguang
 * @Date: 2022-9-26 15:51:00
 * @LastEditors: wenchao
 * @LastEditTime: 2022-9-26 15:51:00
 * @Description:
 */
package switchplatform

type SwitchPlatformReq struct {
	TransId   string `json:"transId"` //切换任务id，管理员绑定端保证不能重复
	NewDomain string `json:"domain"`  //切换目标空间平台域名
}

type SwitchPlatformResp struct {
	TransId    string `json:"transId"`    //切换任务id，请求传入的值。
	UserDomain string `json:"userDomain"` //用户的新域名
}
