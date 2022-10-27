package main

import (
	"flag"
	"log"
	"os"

	"github.com/molu/goc2/pkg/server"
)

func main() {
	LOG_FILE := "./logs/goc2server.log"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("cannot set logfile: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	domainName := flag.String("domain", "localhost", "domain name")
	listenAddr := flag.String("listen", "0.0.0.0", "listen address")
	listenPort := flag.String("port", "4443", "listen port")
	caCertFile := flag.String("cacert", "./certs/ca.crt", "CA certificate")
	serverCertFile := flag.String("cert", "./certs/localhost.crt", "certificate PEM file")
	serverKeyFile := flag.String("key", "./certs/localhost.key", "private key PEM file")
	serverKeyPassFile := flag.String("keypass", "./certs/localhost.pass", "private key passphrase file")

	idHeader := flag.String("ih", "X-Client", "request header containing agent ID")
	cmdHeader := flag.String("ch", "Accept", "response header containing cmd")
	dataHeader := flag.String("dh", "Authorization", "request header containing cmd output")

	s := server.Server{
		DomainName:  *domainName,
		ListenAddr:  *listenAddr + ":" + *listenPort,
		CACertFile:  *caCertFile,
		KeyPassFile: *serverKeyPassFile,
		CertFile:    *serverCertFile,
		PrivKeyFile: *serverKeyFile,
		IDHeader:    *idHeader,
		CmdHeader:   *cmdHeader,
		DataHeader:  *dataHeader,
	}
	s.Run()

}
