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

package certificate

import (
	"agent/biz/model/dto"
	"agent/biz/model/dto/certificate"
	"agent/biz/service/base"
	"agent/config"
	"agent/utils/logger"
	"github.com/dungeonsnd/gocom/encrypt/encoding"
	"io/ioutil"
	//"agent/biz/service/base"
	"bytes"
	cr "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"math/rand"
	"os"
	"time"
)

const (
	certKeyFile = "cert.key"
	certPemFile = "cert.pem"
	certDerFile = "cert.der"
	certCSRFile = "cert.csr"
)

func init() {
	go InitCert()
}

type LanCertService struct {
	base.BaseService
}

type Cert struct {
	Cert           []byte
	CertKey        *rsa.PrivateKey
	CertPem        *bytes.Buffer
	CertKeyPemBuff *bytes.Buffer
	CertKeyPem     []byte
	CertCSR        []byte
	Csr            *x509.Certificate
}

func (svc *LanCertService) Process() dto.BaseRspStr {
	rsp, err := ReadCert(config.Config.Box.Cert.CertDir+certPemFile, config.Config.Box.Cert.CertDir+certKeyFile)
	if err != nil {
		return dto.BaseRspStr{Code: dto.AgentCodeServerErrorStr,
			Message: err.Error()}
	}
	var lancert certificate.LanCert
	lancert.Cert = encoding.Base64Encode(rsp.Raw)
	svc.Rsp = lancert
	return svc.BaseService.Process()
}

var cert Cert

func ReadPriKey() error {
	privateKey, err := os.ReadFile(certKeyFile)
	if err != nil {
		logger.CertificateLogger().Errorf("read private key error:%v", err)
		return err
	}
	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		logger.CertificateLogger().Errorf("failed to decode PEM block containing private key")
		return errors.New("failed to decode PEM block containing private key")
	}
	cert.CertKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	return nil
}

// CreateCSR 创建证书请求文件
func CreateCSR(domains []string) error {
	subj := &pkix.Name{
		Country:            []string{"CN"},
		Organization:       []string{"ISCAS"},
		OrganizationalUnit: []string{"ISRC"},
		Locality:           []string{"Nanjing"},
		Province:           []string{"Jiangsu"},
		StreetAddress:      []string{"ChuangYanRoad"},
		PostalCode:         nil,
		CommonName:         domains[0],
	}
	var err error
	//  生成私钥
	if _, err := os.Stat(config.Config.Box.Cert.CertDir + certKeyFile); err != nil {
		cert.CertKey, err = GeneratePriKey()
		if err != nil {
			logger.CertificateLogger().Errorf("gen rsa error:%+v", err)
			return err
		}
		cert.CertKeyPem = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(cert.CertKey),
		})
		err = ioutil.WriteFile(config.Config.Box.Cert.CertDir+certKeyFile, cert.CertKeyPem, 0644)
		if err != nil {
			logger.CertificateLogger().Errorf("Failed to write KEY to file: %v", err)
			return err
		}
	} else {
		err = ReadPriKey()
		if err != nil {
			logger.CertificateLogger().Errorf("ReadPriKey error:%+v", err)
			return err
		}
	}

	// 生成csr
	csrBytes, err := x509.CreateCertificateRequest(cr.Reader, &x509.CertificateRequest{
		Subject:  *subj,
		DNSNames: domains,
		//SignatureAlgorithm: x509.SHA256WithRSA,
	}, cert.CertKey)
	if err != nil {
		logger.CertificateLogger().Errorf("Failed to create CSR: %v", err)
		return err
	}
	// 写CSR 文件
	cert.CertCSR = pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csrBytes,
	})
	err = ioutil.WriteFile(config.Config.Box.Cert.CertDir+certCSRFile, cert.CertCSR, 0644)
	if err != nil {
		logger.CertificateLogger().Errorf("Failed to write CSR to file: %v", err)
		return err
	}
	//cert.CertKeyPem = new(bytes.Buffer)
	// 写私钥文件

	logger.CertificateLogger().Infof("CSR AND KEY written to file successfully!")
	return nil
}

func InitCert() {

	if _, err := os.Stat(config.Config.Box.Cert.CertDir + certDerFile); err != nil {
		subj := &pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{"ISCAS"},
			OrganizationalUnit: []string{"ISRC"},
			Locality:           []string{"Nanjing"},
			Province:           []string{"Jiangsu"},
			StreetAddress:      []string{"ChuangYanRoad"},
			PostalCode:         nil,
			CommonName:         "AOSPACE LAN CERT",
		}
		// 生成自签名证书
		selfSignCert, err := Req(subj, 10)
		if err != nil {
			logger.CertificateLogger().Errorf("generate self-sign certificate failed")
		}
		err = WriteCert(selfSignCert)
		if err != nil {
			logger.CertificateLogger().Errorf("generate self-sign certificate failed")
		}
	}

}

func GeneratePriKey() (*rsa.PrivateKey, error) {
	prikey, err := rsa.GenerateKey(cr.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return prikey, nil
}

func Req(subj *pkix.Name, expire int) (*Cert, error) {
	var (
		cert = &Cert{}
		err  error
	)
	cert.CertKey, err = GeneratePriKey()
	if err != nil {
		return nil, err
	}
	if expire < 1 {
		expire = 1
	}

	cert.Csr = &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63n(2000)),
		Subject:      *subj,
		//IPAddresses:  ip,
		//DNSNames:     dns,
		IsCA:      false,
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(expire, 0, 0),
		//SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}

	cert.Cert, err = x509.CreateCertificate(cr.Reader, cert.Csr, cert.Csr, &cert.CertKey.PublicKey, cert.CertKey)
	if err != nil {
		return nil, err
	}

	cert.CertPem = new(bytes.Buffer)
	pem.Encode(cert.CertPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Cert,
	})
	cert.CertKeyPemBuff = new(bytes.Buffer)
	pem.Encode(cert.CertKeyPemBuff, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cert.CertKey),
	})

	return cert, nil
}

func ReadCert(certFile, keyFile string) (cert *x509.Certificate, err error) {
	if len(certFile) == 0 && len(keyFile) == 0 {
		return nil, errors.New("cert or key has not provided")
	}
	// load cert and key by tls.LoadX509KeyPair
	tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(tlsCert.Certificate[0])
}

func WriteCert(cert *Cert) error {
	if _, err := os.Stat(config.Config.Box.Cert.CertDir); err != nil {
		err = os.Mkdir(config.Config.Box.Cert.CertDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	kf, err := os.OpenFile(config.Config.Box.Cert.CertDir+certKeyFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer kf.Close()
	if _, err := kf.Write(cert.CertKeyPemBuff.Bytes()); err != nil {
		return err
	}

	cpf, err := os.OpenFile(config.Config.Box.Cert.CertDir+certPemFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer cpf.Close()
	if _, err := cpf.Write(cert.CertPem.Bytes()); err != nil {
		return err
	}

	cf, err := os.OpenFile(config.Config.Box.Cert.CertDir+certDerFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	if _, err := cf.Write(cert.Cert); err != nil {
		return err
	}
	return nil
}
