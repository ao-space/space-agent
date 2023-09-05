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
	"errors"
	"fmt"
	"strings"
)

// TODO: Review and refine comments
// Base prefix according to the specification.
const prefix = "did:"

const didMethod = "aospace"
const verificationMethod = "aospacekey"

const defaultDidFragment = "did0"

var defaultContexts = []interface{}{
	defaultContext,
	securityContext,
}

type Identifier struct {
	data *identifierData
}

type identifierData struct {
	Method string

	ID string

	Params []Param

	Path string

	PathSegments []string

	Query string

	Fragment string

	Controller string

	Context []interface{} `json:"@context" yaml:"-"`

	VerificationMethods []*VerificationKey

	AuthenticationMethod []string

	AssertionMethod []string

	KeyAgreement []string

	CapabilityInvocation []string

	CapabilityDelegation []string
}

func NewIdentifier() (*Identifier, error) {

	d := &identifierData{
		Context:  defaultContexts,
		Fragment: defaultDidFragment,
	}
	return &Identifier{
		data: d,
	}, nil
}

func (d *Identifier) Document(safe bool) *Document {
	doc := &Document{
		Context:              d.data.Context,
		Subject:              d.DID(),
		VerificationMethod:   d.VerificationMethods(),
		Authentication:       d.data.AuthenticationMethod,
		AssertionMethod:      d.data.AssertionMethod,
		KeyAgreement:         d.data.KeyAgreement,
		CapabilityInvocation: d.data.CapabilityInvocation,
		CapabilityDelegation: d.data.CapabilityDelegation,
	}

	return doc
}

func (d *Identifier) VerificationMethods() []*VerificationKey {
	keys := make([]*VerificationKey, len(d.data.VerificationMethods))
	for i, k := range d.data.VerificationMethods {
		keys[i] = k
	}
	return keys
}

func (d *Identifier) DeleteVerificationMethod(index int) error {
	if index < len(d.data.VerificationMethods) {
		d.data.VerificationMethods = append(d.data.VerificationMethods[:index], d.data.VerificationMethods[index+1:]...)
		return nil
	}
	return fmt.Errorf("index:%v error, VerificationMethods len:%v", index, len(d.data.VerificationMethods))
}

func (d *Identifier) DeleteVerificationMethodOfQuery(query string) (int, error) {
	if d.data == nil || d.data.VerificationMethods == nil {
		return 0, nil
	}

	vks := make([]*VerificationKey, 0, len(d.data.VerificationMethods))
	for _, k := range d.data.VerificationMethods {
		if strings.Index(k.ID, query) < 0 {
			vks = append(vks, k)
		}
	}
	cnt := len(d.data.VerificationMethods) - len(vks)
	d.data.VerificationMethods = vks
	return cnt, nil
}

func (d *Identifier) GetVerificationMethodCountOfPublicKey() int {
	cnt := 0
	for _, k := range d.data.VerificationMethods {
		if k.Type != KeyTypeMutisig {
			cnt++
		}
	}
	return cnt
}

func (d *Identifier) AddNewCapabilityInvocation() error {

	d.data.CapabilityInvocation = append(d.data.CapabilityInvocation, "#multisig-0")
	d.update()
	return nil
}

func (d *Identifier) AddNewVerificationMethod(keyType, publicKey, query, fragment string) (string, error) {
	return d.AddNewVerificationMethodWithIndex(-1, keyType, publicKey, query, fragment)
}

func (d *Identifier) RemoveVerificationMethodByFragment(fragment string) (int, error) {
	if d.data == nil || d.data.VerificationMethods == nil {
		return 0, nil
	}

	vks := make([]*VerificationKey, 0, len(d.data.VerificationMethods))
	for _, k := range d.data.VerificationMethods {
		if k.Fragment() != fragment {
			vks = append(vks, k)
		}
	}
	cnt := len(d.data.VerificationMethods) - len(vks)
	d.data.VerificationMethods = vks
	return cnt, nil
}

func (d *Identifier) AddNewVerificationMethodWithIndex(index int, keyType, publicKey, query, fragment string) (string, error) {

	id := d.VerificationDID(CalVerificationIdString(publicKey))
	for _, k := range d.data.VerificationMethods {
		if k.IdString() == id {
			return "", errors.New("duplicated key identifier")
		}
	}
	pk, err := newVerificationMethodByPublicKey(id, keyType, publicKey, query, fragment)
	if err != nil {
		return "", err
	}
	pk.Controller = d.Fragment()
	if index < 0 || index >= len(d.data.VerificationMethods) {
		d.data.VerificationMethods = append(d.data.VerificationMethods, pk)
	} else {
		d.data.VerificationMethods = append(d.data.VerificationMethods, pk)
		copy(d.data.VerificationMethods[index+1:], d.data.VerificationMethods[index:])
		d.data.VerificationMethods[index] = pk
	}
	d.update()
	return id, nil
}

func (d *Identifier) AddNewVerificationMethodOfMultisig(firstVerificationMethod, secondVerificationMethod []string) error {

	controller := d.Fragment()
	pk, err := newVerificationMethodOfMultisig(controller, firstVerificationMethod, secondVerificationMethod)
	if err != nil {
		return err
	}
	d.data.VerificationMethods = append(d.data.VerificationMethods, pk)
	d.update()
	return nil
}

func (d *Identifier) DID() string {
	return fmt.Sprintf("%s%s:%s%s", prefix, didMethod, d.idString(), d.Fragment())
}

func (d *Identifier) Fragment() string {
	if d.data.Fragment == "" {
		return ""
	}
	return fmt.Sprintf("#%s", d.data.Fragment)
}

func (d *Identifier) String() string {
	var buf strings.Builder

	buf.WriteString(d.DID())

	if d.data.Query != "" {
		// write a leading ? and then Query
		buf.WriteByte('?')
		buf.WriteString(d.data.Query)
	}

	if d.data.Fragment != "" {
		// write a leading # and then the fragment value
		buf.WriteByte('#')
		buf.WriteString(d.data.Fragment)
	}

	return buf.String()
}

func (d *Identifier) VerificationDID(idString string) string {
	return fmt.Sprintf("%s%s:%s", prefix, verificationMethod, idString)
}

func (d *Identifier) idString() string {
	return d.data.ID
}

func (d *Identifier) params() string {
	if len(d.data.Params) == 0 {
		return ""
	}

	var buf strings.Builder
	for _, p := range d.data.Params {
		if param := p.String(); param != "" {
			buf.WriteByte(';')
			buf.WriteString(param)
		}
	}
	return buf.String()
}

func (d *Identifier) update() {
	methods := d.VerificationMethods()
	for _, m := range methods {
		// TODO: 需要判断改密钥是否在 CapabilityInvocation 的定义中。由于时间关系，暂时取第一个公钥。后续需要完善 !
		if len(m.PublicKeyPem) > 0 {
			d.data.ID = CalAOSpaceIdString(m.PublicKeyPem) // use first one temp.
			break
		}
	}

}

func FromDocument(doc *Document) (*Identifier, error) {
	id, err := Parse(doc.Subject)
	if err != nil {
		return nil, err
	}

	// Restore public keys
	for _, k := range doc.VerificationMethod {
		rk := &VerificationKey{}
		rk = k
		id.data.VerificationMethods = append(id.data.VerificationMethods, rk)
	}

	// Restore verification relationships
	id.data.Context = doc.Context
	id.data.Controller = doc.Controller
	id.data.AuthenticationMethod = append(id.data.AuthenticationMethod, doc.Authentication...)
	id.data.AssertionMethod = append(id.data.AssertionMethod, doc.AssertionMethod...)
	id.data.KeyAgreement = append(id.data.KeyAgreement, doc.KeyAgreement...)
	id.data.CapabilityInvocation = append(id.data.CapabilityInvocation, doc.CapabilityInvocation...)
	id.data.CapabilityDelegation = append(id.data.CapabilityDelegation, doc.CapabilityDelegation...)
	return id, nil
}
