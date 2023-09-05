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

package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 生成JWT
// expiresAt : time.Now().Add(24 * time.Hour * 365 * 5)
func GenerateJWT(issuer, subject string, audience []string, expiresAt time.Time, mapPrivate, mapPublic map[string]string, key *rsa.PrivateKey) (string, error) {
	claims := jwt.MapClaims{
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": issuer,
		"sub": subject,
		"aud": audience,
	}
	for k, v := range mapPrivate {
		claims[k] = v
	}
	for k, v := range mapPublic {
		claims[k] = v
	}

	if key == nil {
		// 使用HS256签名算法
		token := jwt.NewWithClaims(SigningMethodCustom, claims)
		// token.SigningString() this is the data to sign.
		signedJWT, err := token.SignedString(key)
		return signedJWT, err
	} else {
		// 使用HS256签名算法
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		// token.SigningString() this is the data to sign.
		signedJWT, err := token.SignedString(key)
		return signedJWT, err
	}
}

// 解析JWT
// return issuer, subject, audience, map, error
func ParseJwt(tokenString string, publicKey *rsa.PublicKey) (string, string, []string, map[string]string, error) {
	// t, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return publicKey, nil
	// }, jwt.WithLeeway(5*time.Second))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if publicKey == nil {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*SigningMethodCustomSt); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
		} else {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return publicKey, nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return "", "", nil, nil, err
	}
	if token == nil {
		return "", "", nil, nil, fmt.Errorf("parse returned token is nil")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		issuer, err := claims.GetIssuer()
		if err != nil {
			return "", "", nil, nil, err
		}
		subject, err := claims.GetSubject()
		if err != nil {
			return "", "", nil, nil, err
		}
		audience, err := claims.GetAudience()
		if err != nil {
			return "", "", nil, nil, err
		}

		m := make(map[string]string)
		for k := range claims {
			if k == "exp" || k == "nbf" || k == "iat" || k == "aud" || k == "iss" || k == "sub" {
				continue
			}
			m[k] = claims[k].(string)
		}
		return issuer, subject, audience, m, nil
	} else {
		return "", "", nil, nil, err
	}

	// if claims, ok := t.Claims.(*MyCustomClaims); ok && t.Valid {
	// 	return claims.EncodeObjPrivate, claims.EncodeObjPublic, claims.Issuer, claims.Subject, claims.Audience, nil
	// } else {
	// 	return "", "", "", "", nil, err
	// }
}
