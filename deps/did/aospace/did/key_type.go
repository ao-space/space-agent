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
	"encoding/json"
	"fmt"

	"github.com/mr-tron/base58"
)

type KeyType int

const (
	KeyTypeEd KeyType = iota

	KeyTypeRSA

	KeyTypeSecp256k1

	KeyTypeMutisig
)

func (v KeyType) String() string {
	values := [...]string{
		"Ed25519VerificationKey2020",
		"RsaVerificationKey2018",
		"EcdsaSecp256k1VerificationKey2019",
		"ConditionalProof2022",
	}
	if int(v) > len(values) {
		return "unknown key type"
	}
	return values[v]
}

func (v KeyType) SignatureType() string {
	values := [...]string{
		"Ed25519Signature2020",
		"RsaSignature2018",
		"EcdsaSecp256k1Signature2019",
	}
	if int(v) > len(values) {
		return "unknown signature type"
	}
	return values[v]
}

func (v KeyType) EncodePublicKey(vk *VerificationKey, pub []byte) {
	if v == KeyTypeEd {
		vk.Public = multibaseEncode(pub)
		return
	}
	vk.PublicKeyPem = string(pub)
	// vk.PublicKeyBase58 = base58.Encode(pub)
}

func (v KeyType) DecodePublicKey(vk *VerificationKey) ([]byte, error) {
	if v == KeyTypeEd {
		return multibaseDecode(vk.Public)
	}
	return base58.Decode(vk.PublicKeyBase58)
}

func (v *KeyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *KeyType) UnmarshalJSON(b []byte) error {
	var s string
	var err error
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*v, err = keyTypeFromString(s)
	if err != nil {
		return err
	}
	return nil
}

func (v KeyType) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

func (v *KeyType) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var s string
	var err error
	if err = unmarshal(&s); err != nil {
		return err
	}
	*v, err = keyTypeFromString(s)
	if err != nil {
		return err
	}
	return nil
}

func keyTypeFromString(val string) (kt KeyType, err error) {
	switch val {
	case KeyTypeEd.String():
		kt = KeyTypeEd
		return
	case KeyTypeRSA.String():
		kt = KeyTypeRSA
		return
	case KeyTypeSecp256k1.String():
		kt = KeyTypeSecp256k1
		return
	case KeyTypeMutisig.String():
		kt = KeyTypeMutisig
		return
	default:
		err = fmt.Errorf("unknown key type: %s", val)
		return
	}
}
