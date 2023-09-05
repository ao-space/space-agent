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

import "testing"

func TestCalAOSpaceIdString(t *testing.T) {

	publicKey := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnN5jap7CGcqYURbLDVUa\nLc9kMxOyCMEykfwbQKXvTkPMkR9tKZmq8EqfG2d2OyUpF1TIfqHK7Q6d33yD02oO\nBTXZw1Ijkfxvu0KwG2zLV02FTuwZzgYa/AaP5iRZDx5GwTk/YFw+NTqT8Gf29a/L\n/ItcCfsEFLr3zMDXUcU9A7rBEy5ncva6RLNpXawegFGlCZa5+Gah8voKl8ZGpIgt\nlSc1IdnbPbBCYYlUATWLCLeYl+Q9/LslbpkFtdR+4M8vU7G1H+AQZ5fr2E9qX36I\nzcnchDmKq5bkbWQ9GJeZKqZTkhtCPBy4cphM8fHtZuoh1fA3VfF01N4KHT2bUdtp\nJwIDAQAB\n-----END PUBLIC KEY-----"

	idstring := CalAOSpaceIdString(publicKey)
	t.Logf("idstring:\n%+v\n\n", idstring)

}

func TestCalVerificationIdString(t *testing.T) {

	publicKey := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnN5jap7CGcqYURbLDVUa\nLc9kMxOyCMEykfwbQKXvTkPMkR9tKZmq8EqfG2d2OyUpF1TIfqHK7Q6d33yD02oO\nBTXZw1Ijkfxvu0KwG2zLV02FTuwZzgYa/AaP5iRZDx5GwTk/YFw+NTqT8Gf29a/L\n/ItcCfsEFLr3zMDXUcU9A7rBEy5ncva6RLNpXawegFGlCZa5+Gah8voKl8ZGpIgt\nlSc1IdnbPbBCYYlUATWLCLeYl+Q9/LslbpkFtdR+4M8vU7G1H+AQZ5fr2E9qX36I\nzcnchDmKq5bkbWQ9GJeZKqZTkhtCPBy4cphM8fHtZuoh1fA3VfF01N4KHT2bUdtp\nJwIDAQAB\n-----END PUBLIC KEY-----"

	idstring := CalVerificationIdString(publicKey)
	t.Logf("Verification idstring:\n%+v\n\n", idstring)

}
