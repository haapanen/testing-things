package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

type KeyGenerator struct {
	caKey  ed25519.PrivateKey
	caCert *x509.Certificate
}

// constructor
func NewKeyGenerator(caCertPath string, caKeyPath string) *KeyGenerator {
	caKeyBytes, err := os.ReadFile(caKeyPath)
	if err != nil {
		panic(err)
	}

	caCertBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		panic(err)
	}

	caKeyBlock, _ := pem.Decode(caKeyBytes)
	if caKeyBlock == nil {
		panic("failed to parse PEM block containing the key")
	}

	caCertBlock, _ := pem.Decode(caCertBytes)
	if caCertBlock == nil {
		panic("failed to parse PEM block containing the certificate")
	}

	caKey, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		panic(err)
	}

	return &KeyGenerator{
		caKey:  caKey.(ed25519.PrivateKey),
		caCert: caCert,
	}
}

type KeyPair struct {
	PrivateKey   string  `json:"privateKey"`
	PublicKey    string  `json:"publicKey"`
	SerialNumber big.Int `json:"serialNumber"`
}

func (k *KeyGenerator) GenerateKeysAndCertificates() (keyPair *KeyPair, err error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	clientKeyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}

	clientKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: clientKeyBytes,
	})

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "IoT Device",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	clientCertBytes, err := x509.CreateCertificate(rand.Reader, template, k.caCert, pub, k.caKey)
	if err != nil {
		return nil, err
	}

	clientCert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCertBytes,
	})

	return &KeyPair{
		// convert to string
		PrivateKey:   string(clientKey),
		PublicKey:    string(clientCert),
		SerialNumber: *serialNumber,
	}, nil
}
