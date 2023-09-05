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
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/mr-tron/base58"
)

func getHash(data []byte) []byte {
	h := sha256.New()
	if _, err := h.Write(data); err != nil {
		return nil
	}
	return h.Sum(nil)
}

func multibaseEncode(data []byte) string {
	return "z" + base58.Encode(data)
}

func multibaseDecode(src string) ([]byte, error) {
	base := src[:1]
	data := src[1:]
	switch base {
	case "z":
		return base58.Decode(data)
	case "f":
		return hex.DecodeString(data)
	case "m":
		return base64.RawStdEncoding.DecodeString(data)
	case "M":
		return base64.StdEncoding.DecodeString(data)
	case "u":
		return base64.RawURLEncoding.DecodeString(data)
	case "U":
		return base64.URLEncoding.DecodeString(data)
	default:
		return nil, fmt.Errorf("unsupported base identifier: %s", base)
	}
}
