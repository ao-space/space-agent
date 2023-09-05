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
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

const (
	defaultMultisigFragment = "multisig-0"
)

type VerificationKey struct {
	ID string `json:"id" yaml:"id"`

	Type KeyType `json:"type" yaml:"type"`

	Controller string `json:"controller" yaml:"controller"`

	Public string `json:"publicKeyMultibase,omitempty" yaml:"publicKeyMultibase,omitempty"`

	PublicKeyPem string `json:"publicKeyPem,omitempty" yaml:"publicKeyPem,omitempty"`

	PublicKeyBase58 string `json:"publicKeyBase58,omitempty" yaml:"publicKeyBase58,omitempty"`

	ConditionOr  []interface{} `json:"conditionOr,omitempty" yaml:"conditionOr,omitempty"`
	ConditionAnd []interface{} `json:"conditionAnd,omitempty" yaml:"conditionAnd,omitempty"`
}

func (vk *VerificationKey) IdString() string {
	if i := strings.Index(vk.ID, "#"); i > 0 {
		return vk.ID[:i]
	}
	return ""
}
func (vk *VerificationKey) Fragment() string {
	if i := strings.Index(vk.ID, "#"); i > 0 {
		return vk.ID[i:]
	}
	return ""
}

func newVerificationMethodByPublicKey(id, keyType, publicKeyPem, query, fragment string) (*VerificationKey, error) {

	var kt KeyType
	if KeyTypeRSA.String() == keyType {
		kt = KeyTypeRSA
	} else {
		return nil, fmt.Errorf("unsupprted keytype of %v", keyType)
	}

	vk := &VerificationKey{Type: kt, ID: id + "?" + query + "#" + fragment}
	switch kt {
	case KeyTypeRSA:
		kt.EncodePublicKey(vk, []byte(publicKeyPem))
	default:
		return nil, fmt.Errorf("unsupprted keytype of %v", kt)
	}

	return vk, nil
}

func getArrayExcludeElement(firstVerificationMethod, secondVerificationMethod []string, exclude string) []interface{} {
	res := make([]interface{}, 0)
	for _, v := range firstVerificationMethod {
		if v != exclude {
			res = append(res, v)
		}
	}
	for _, v := range secondVerificationMethod {
		if v != exclude {
			res = append(res, v)
		}
	}
	return res
}

func newVerificationMethodOfMultisig(controller string, firstVerificationMethod, secondVerificationMethod []string) (*VerificationKey, error) {
	if len(firstVerificationMethod) < 1 && len(firstVerificationMethod)+len(secondVerificationMethod) < 2 {
		return nil, fmt.Errorf("VerificationMethod size error. len(firstVerificationMethod):%v, len(secondVerificationMethod):%v",
			len(firstVerificationMethod), len(secondVerificationMethod))
	}

	// 比较灵活的组装方式.
	vkPrimaries := make([]interface{}, 0)
	for _, primary := range firstVerificationMethod {
		others := getArrayExcludeElement(firstVerificationMethod, secondVerificationMethod, primary)
		vkOthers := &VerificationKey{ID: CalVerificationIdString(uuid.NewV4().String()),
			Type:        KeyTypeMutisig,
			Controller:  controller,
			ConditionOr: others} // e.g. key-1 || key-2

		vkPrimary := &VerificationKey{ID: CalVerificationIdString(uuid.NewV4().String()),
			Type:         KeyTypeMutisig,
			Controller:   controller,
			ConditionAnd: []interface{}{primary, vkOthers}} // e.g. key-0 && (key-1 || key-2)
		vkPrimaries = append(vkPrimaries, vkPrimary)
	}

	// 合并
	// e.g. (key-0 && (key-1 || key-2)) || (key-1 && (key-0 || key-2))
	vk := &VerificationKey{ID: CalVerificationIdString(uuid.NewV4().String()) + "#" + defaultMultisigFragment,
		Type:        KeyTypeMutisig,
		Controller:  controller,
		ConditionOr: vkPrimaries}

	return vk, nil
}
