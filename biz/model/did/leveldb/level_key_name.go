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

package leveldb

const (
	// 大类
	prefixPriKey = "key"
	prefixDidDoc = "doc"
	prefixIndex  = "index"
)

func KNameOfSpaceRSAPri(aoId string) string {
	return prefixPriKey + "--space_rsa_pri--" + aoId
}
func KNameOfPasswordRSAPri(aoId string) string {
	return prefixPriKey + "--password_rsa_pri--" + aoId
}
func KNameOfDidDoc(did string) string {
	return prefixDidDoc + "--did_doc--" + did
}

func KNameOfAoIdToDid(aoId string) string {
	return prefixIndex + "--aoid_to_did--" + aoId
}
func KNameOfDidToAoId(did string) string {
	return prefixIndex + "--did_to_aoid--" + did
}