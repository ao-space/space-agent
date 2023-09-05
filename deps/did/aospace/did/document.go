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

type Document struct {
	Context []interface{} `json:"@context" yaml:"-"`

	Subject string `json:"id" yaml:"id"`

	Controller string `json:"controller,omitempty" yaml:"controller,omitempty"`

	AlsoKnownAs []string `json:"alsoKnownAs,omitempty" yaml:"alsoKnownAs,omitempty"`

	VerificationMethod []*VerificationKey `json:"verificationMethod,omitempty" yaml:"verificationMethod,omitempty"`

	Authentication []string `json:"authentication,omitempty" yaml:"authentication,omitempty"`

	AssertionMethod []string `json:"assertionMethod,omitempty" yaml:"assertionMethod,omitempty"`

	KeyAgreement []string `json:"keyAgreement,omitempty" yaml:"keyAgreement,omitempty"`

	CapabilityInvocation []string `json:"capabilityInvocation,omitempty" yaml:"capabilityInvocation,omitempty"`

	CapabilityDelegation []string `json:"capabilityDelegation,omitempty" yaml:"capabilityDelegation,omitempty"`
}

type DocumentMetadata struct {
	Created string `json:"created,omitempty" yaml:"created,omitempty"`

	Updated string `json:"updated,omitempty" yaml:"updated,omitempty"`

	Deactivated bool `json:"deactivated" yaml:"deactivated"`
}

func (d *Document) ExpandedLD() ([]byte, error) {
	return expand(d)
}

func (d *Document) NormalizedLD() ([]byte, error) {
	return normalize(d)
}
