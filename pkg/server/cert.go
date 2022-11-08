package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func (s *Server) getTLSConfig() (config *tls.Config) {
	certPEMBlock, err := os.ReadFile(s.CertFile)
	if err != nil {
		log.Fatalf("cannot read the server cert file: %v", err)
	}

	encryptedKeyBlock, err := os.ReadFile(s.PrivKeyFile)
	if err != nil {
		log.Fatalf("cannot read the server key file: %v", err)
	}

	passphrase, err := os.ReadFile(s.KeyPassFile)
	if err != nil {
		log.Fatalf("cannot read the server key passphrase file: %v", err)
	}

	privateKeyBlock := parseECPrivateKeyWithPasphrase(encryptedKeyBlock, passphrase)
	serverCert, err := tls.X509KeyPair(certPEMBlock, privateKeyBlock)
	if err != nil {
		log.Fatalf("cannot get the server cert: %v", err)
	}

	config = &tls.Config{
		ServerName:               s.DomainName,
		Certificates:             []tls.Certificate{serverCert},
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                s.getCertPool(),
		MinVersion:               tls.VersionTLS13,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	return
}

func (s *Server) getCertPool() (certPool *x509.CertPool) {
	certPool = x509.NewCertPool()

	caCert, err := os.ReadFile(s.CACertFile)
	if err != nil {
		log.Fatalf("failed to load CA cert file: %v", err)
	}
	certPool.AppendCertsFromPEM(caCert)
	return
}

func parseECPrivateKeyWithPasphrase(encryptedPrivateKeyBytes, password []byte) (privateKeyBytes []byte) {
	block, _ := pem.Decode(encryptedPrivateKeyBytes)
	if block == nil {
		log.Fatalf("failed to decode private key to PEM")
	}

	decryptedPrivateKeyBytes, err := x509.DecryptPEMBlock(block, password)
	if err != nil {
		log.Fatalf("failed to decrypt private key: %v", err)
	}

	ecdsaPrivateKey, err := x509.ParseECPrivateKey(decryptedPrivateKeyBytes)
	if err != nil {
		log.Fatalf("failed to parse PEM block EC private key: %v", err)
	}

	privateKeyDERBytes, err := x509.MarshalECPrivateKey(ecdsaPrivateKey)
	if err != nil {
		log.Fatalf("failed to marshal EC private key to DER: %v", err)
	}

	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyDERBytes,
	}

	privateKeyBytes = pem.EncodeToMemory(keyBlock)

	return
}
