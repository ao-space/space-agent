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

package jwt

import (
	utilshttp "agent/utils/network/http"
	"fmt"
)

func PostHttpRequest(url string,
	parms interface{},
	headers map[string]string,
	res interface{}) error {
	return SendHttpRequest("POST", url, parms, headers, res)
}

func GetHttpRequest(url string,
	parms interface{},
	headers map[string]string,
	res interface{}) error {
	return SendHttpRequest("GET", url, parms, headers, res)
}

func SendHttpRequest(method string, url string,
	parms interface{},
	headers map[string]string,
	res interface{}) error {

	if parms == nil {
		parms = make(map[string]string)
	}

	req, rsp, content, err := utilshttp.SendJsonWithHeaders(method, url, parms, headers, res)
	if err != nil {
		return fmt.Errorf("http.SendHttpRequest Failed, err=%v. INPUT [%v] url=%+v, parms=%+v, headers=%+v. OUTPUT req=%+v, rsp=%+v, content=%+v",
			err, method, url, parms, headers, req, rsp, string(content))
	}

	return nil
}
