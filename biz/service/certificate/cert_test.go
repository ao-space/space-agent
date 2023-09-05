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

package certificate

import (
	"fmt"
	"regexp"
	"testing"
)

func TestDNSTXTReg(t *testing.T) {
	str := "[Thu Mar 23 11:36:05 CST 2023] Getting domain auth token for each domain\n[Thu Mar 23 11:36:05 CST 2023] Getting webroot for domain='lan.sit-space.eulix.xyz'\n[Thu Mar 23 11:36:05 CST 2023] Add the following TXT record:\n[Thu Mar 23 11:36:05 CST 2023] Domain: '_acme-challenge.xxxx.sit-space.eulix.xyz'\n[Thu Mar 23 11:36:05 CST 2023] TXT value: '9B5oxxxxxxxxxxxxxxxxxxxxxxxxx'\n[Thu Mar 23 11:36:05 CST 2023] Please be aware that you prepend _acme-challenge. before your domain\n[Thu Mar 23 11:36:05 CST 2023] so the resulting subdomain will be: _acme-challenge.lan.sit-space.eulix.xyz\n[Thu Mar 23 11:36:05 CST 2023] Please add the TXT records to the domains, and re-run with --renew.\n[Thu Mar 23 11:36:05 CST 2023] Please add '--debug' or '--log' to check more details."

	// 正则表达式，用于匹配 Domain 和 TXT value
	re := regexp.MustCompile(`Domain:\s*'([^']*)'\n.*TXT value: *'([^']*)'`)

	// 查找第一个匹配的结果
	match := re.FindStringSubmatch(str)

	if len(match) > 0 {
		// 输出匹配结果
		fmt.Printf("Domain: %s\n", match[1])
		fmt.Printf("TXT value: %s\n", match[2])
	} else {
		fmt.Println("No match found")
	}

}

func TestCheckDNSTXT(t *testing.T) {
	CheckDNSTXT("xuyangtest.lan.sit-space.eulix.xyz")
}
