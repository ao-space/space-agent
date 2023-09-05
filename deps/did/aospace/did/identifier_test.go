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
	"testing"
)

func TestIdentifier(t *testing.T) {
	did, err := NewIdentifier()
	if err != nil {
		panic(err)
	}
	t.Logf("did:%+v", did)

	keyType := "RsaVerificationKey2018"
	publicKeyPem := "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnN5jap7CGcqYURbLDVUa\nLc9kMxOyCMEykfwbQKXvTkPMkR9tKZmq8EqfG2d2OyUpF1TIfqHK7Q6d33yD02oO\nBTXZw1Ijkfxvu0KwG2zLV02FTuwZzgYa/AaP5iRZDx5GwTk/YFw+NTqT8Gf29a/L\n/ItcCfsEFLr3zMDXUcU9A7rBEy5ncva6RLNpXawegFGlCZa5+Gah8voKl8ZGpIgt\nlSc1IdnbPbBCYYlUATWLCLeYl+Q9/LslbpkFtdR+4M8vU7G1H+AQZ5fr2E9qX36I\nzcnchDmKq5bkbWQ9GJeZKqZTkhtCPBy4cphM8fHtZuoh1fA3VfF01N4KHT2bUdtp\nJwIDAQAB\n-----END PUBLIC KEY-----"

	did.AddNewVerificationMethod(keyType, publicKeyPem, "k=v", "key-0")
	did.AddNewVerificationMethodOfMultisig()

	didDoc := did.Document(true)
	js, err := json.Marshal(didDoc)
	if err != nil {
		panic(err)
	}
	t.Logf("didDoc json Marshal :%+v", string(js))

	normalDidDoc, err := didDoc.NormalizedLD()
	if err != nil {
		panic(err)
	}
	t.Logf("normalDidDoc:%+v", string(normalDidDoc))
}

func TestFromDocument(t *testing.T) {
	doc := `{
		"@context": [
			"https://www.w3.org/ns/did/v1",
			"https://w3id.org/security/v1"
		],
		"id": "did:aospace:11MSzi1eoQtDfpw6ez6AbkaPWvQBaR2tuki#did0",
	
		"verificationMethod": [{
				"id": "did:aospacekey:AAAHtMWCPnvz2q5ONvw=?versionTime=2021-05-10T17:00:00Z&credentialType=device#key-0",
				"type": "RsaVerificationKey2018",
				"controller": "#did0",
				"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnN5jap7CGcqYURbLDVUa\nLc9kMxOyCMEykfwbQKXvTkPMkR9tKZmq8EqfG2d2OyUpF1TIfqHK7Q6d33yD02oO\nBTXZw1Ijkfxvu0KwG2zLV02FTuwZzgYa/AaP5iRZDx5GwTk/YFw+NTqT8Gf29a/L\n/ItcCfsEFLr3zMDXUcU9A7rBEy5ncva6RLNpXawegFGlCZa5+Gah8voKl8ZGpIgt\nlSc1IdnbPbBCYYlUATWLCLeYl+Q9/LslbpkFtdR+4M8vU7G1H+AQZ5fr2E9qX36I\nzcnchDmKq5bkbWQ9GJeZKqZTkhtCPBy4cphM8fHtZuoh1fA3VfF01N4KHT2bUdtp\nJwIDAQAB\n-----END PUBLIC KEY-----"
			},
			{
				"id": "did:aospacekey:AAAHtMWCPnvz2q5ONvw=?versionTime=2021-05-10T17:00:00Z&credentialType=binder&deviceName=iPhoneX#key-1",
				"type": "RsaVerificationKey2018",
				"controller": "#did0",
				"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnN5jap7CGcqYURbLDVUa\nLc9kMxOyCMEykfwbQKXvTkPMkR9tKZmq8EqfG2d2OyUpF1TIfqHK7Q6d33yD02oO\nBTXZw1Ijkfxvu0KwG2zLV02FTuwZzgYa/AaP5iRZDx5GwTk/YFw+NTqT8Gf29a/L\n/ItcCfsEFLr3zMDXUcU9A7rBEy5ncva6RLNpXawegFGlCZa5+Gah8voKl8ZGpIgt\nlSc1IdnbPbBCYYlUATWLCLeYl+Q9/LslbpkFtdR+4M8vU7G1H+AQZ5fr2E9qX36I\nzcnchDmKq5bkbWQ9GJeZKqZTkhtCPBy4cphM8fHtZuoh1fA3VfF01N4KHT2bUdtp\nJwIDAQAB\n-----END PUBLIC KEY-----"
			},
			{
				"id": "did:aospacekey:AAAHtMWCPnvz2q5ONvw=?versionTime=2021-05-10T17:00:00Z&credentialType=password#key-2",
				"type": "RsaVerificationKey2018",
				"controller": "#did0",
				"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnN5jap7CGcqYURbLDVUa\nLc9kMxOyCMEykfwbQKXvTkPMkR9tKZmq8EqfG2d2OyUpF1TIfqHK7Q6d33yD02oO\nBTXZw1Ijkfxvu0KwG2zLV02FTuwZzgYa/AaP5iRZDx5GwTk/YFw+NTqT8Gf29a/L\n/ItcCfsEFLr3zMDXUcU9A7rBEy5ncva6RLNpXawegFGlCZa5+Gah8voKl8ZGpIgt\nlSc1IdnbPbBCYYlUATWLCLeYl+Q9/LslbpkFtdR+4M8vU7G1H+AQZ5fr2E9qX36I\nzcnchDmKq5bkbWQ9GJeZKqZTkhtCPBy4cphM8fHtZuoh1fA3VfF01N4KHT2bUdtp\nJwIDAQAB\n-----END PUBLIC KEY-----"
			},
			{
				"id": "did:aospacekey:uuid?versionTime=2021-05-10T17:00:00Z#multisig-0",
				"controller": "#did0",
				"type": "ConditionalProof2022",
				"conditionOr": [{
						"id": "did:aospacekey:uuid?versionTime=2021-05-10T17:00:00Z#multisig-0-0",
						"controller": "#did0",
						"type": "ConditionalProof2022",
						"conditionAnd": [
							"#key-0",
							{
								"id": "did:aospacekey:uuid?versionTime=2021-05-10T17:00:00Z#multisig-0-0-0",
								"controller": "#did0",
								"type": "ConditionalProof2022",
								"conditionOr": [
									"#key-1", "#key-2"
								]
							}
						]
					},
					{
						"id": "did:aospacekey:uuid?versionTime=2021-05-10T17:00:00Z#multisig-0-1",
						"controller": "#did0",
						"type": "ConditionalProof2022",
						"conditionAnd": [
							"#key-1",
							{
								"id": "did:aospacekey:uuid?versionTime=2021-05-10T17:00:00Z#multisig-0-1-0",
								"controller": "#did0",
								"type": "ConditionalProof2022",
								"conditionOr": [
									"#key-0", "#key-2"
								]
							}
						]
					}
				]
			}
		],
	
		"capabilityInvocation": [
			"#multisig-0"
		]
	}
	`

	didDoc := &Document{}
	err := json.Unmarshal([]byte(doc), didDoc)
	if err != nil {
		panic(err)
	}
	t.Logf("didDoc:%+v", didDoc)

	didObj, err := FromDocument(didDoc)
	if err != nil {
		panic(err)
	}
	t.Logf("didObj.DID:%+v", didObj.DID())

	methods := didObj.VerificationMethods()
	for _, m := range methods {
		t.Logf("m: %+v", m)
	}

	if len(methods) > 1 {
		err = didObj.DeleteVerificationMethodOfQuery("credentialType=password")
		if err != nil {
			panic(err)
		}
	}

	keyType := "RsaVerificationKey2018"
	publicKeyPem := "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDV1/eYapPImz/PZ47cL5hzNeWl\nDoTi+x3bfwjV7Q0NJldyzGkayRH3Jdtwpkm47Txicr767O22Y7hzE9OrobP+zx9m\n7s86Xd6kHbj+mSzySKWr7AP7/RANIGns/rDh5bMzPdbtcD9LTn2SLgcAE3RSTKY3\n4noqfedxX/ko3sQXYwIDAQAB\n-----END PUBLIC KEY-----"
	newKeyId, err := didObj.AddNewVerificationMethodWithIndex(didObj.GetVerificationMethodCountOfPublicKey(),
		keyType, publicKeyPem, "k=v", "key-1-new")
	if err != nil {
		panic(err)
	}
	t.Logf("\n\n==== newKeyId:%+v\n\n", newKeyId)

	methods = didObj.VerificationMethods()
	for _, m := range methods {
		t.Logf("m: %+v", m)
	}

	newDidDoc := didObj.Document(true)
	newDidDocBytes, err := json.Marshal(newDidDoc)
	if err != nil {
		panic(err)
	}
	t.Logf("\n\nnewDidDocBytes :%+v\n\n", string(newDidDocBytes))
}
