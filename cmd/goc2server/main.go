package main

import (
	"flag"

	"github.com/molu/goc2/pkg/server"
)

func main() {
	listenAddr := flag.String("listen", "0.0.0.0", "listen address")
	listenPort := flag.String("port", "4443", "listen port")
	domainName := flag.String("domain", "localhost", "domain name")

	serverCertFile := flag.String("cert", "./certs/localhost.crt", "certificate PEM file")
	serverKeyFile := flag.String("key", "./certs/localhost.key", "private key PEM file")
	serverKeyPassFile := flag.String("keypass", "./certs/localhost.pass", "private key passphrase file")
	caCertFile := flag.String("cacert", "./certs/ca.crt", "CA certificate PEM file")

	idHeader := flag.String("ih", "X-Client", "request header containing agent ID")
	cmdHeader := flag.String("ch", "Accept", "response header containing command")
	dataHeader := flag.String("dh", "Authorization", "request header containing command output")

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
	s.Listen()
}
