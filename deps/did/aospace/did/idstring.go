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

package did

import (
	"encoding/base64"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

func CalAOSpaceIdString(publicKey string) string {

	hash_sha256 := sha3.Sum256([]byte(publicKey))

	hasher := ripemd160.New()
	hasher.Write(hash_sha256[0:])
	hash_ripemd160 := hasher.Sum(nil)

	version := []byte{0, 0}

	data := version
	data = append(data, hash_ripemd160...)

	digest := sha3.Sum256(data)
	checksum := digest[0:4]

	data = append(data, checksum...)
	idstring := base58.Encode(data)

	return idstring
}

func CalVerificationIdString(content string) string {
	hash_sha256 := sha3.Sum256([]byte(content))
	hash_sha256_cutted := hash_sha256[0:8]

	version := []byte{0, 0}

	data := version
	data = append(data, hash_sha256_cutted...)

	digest := sha3.Sum256(data)
	checksum := digest[0:4]

	data = append(data, checksum...)
	idstring := base64.StdEncoding.EncodeToString(data)

	return idstring
}
