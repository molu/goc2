#!/bin/bash
# goc2server listen address
SERVER_ADDR=https://localhost:4443
# certificates and private keys dir
CERTS_DIR=certs
# agent unique name (e.g. UUIDv4)
AGENT_NAME=testing  # $(uuidgen -r)
# request header containing agent name
ID_HEADER=X-Client
# response header containing cmd to execute
CMD_HEADER=Accept
# request header to send cmd output in
DATA_HEADER=Authorization
# delay between requests in miliseconds
REQUEST_DELAY=10000  

bash ${PWD}/bin/certgen.sh ${AGENT_NAME}

CA_CERT_B64=$(cat ${CERTS_DIR}/ca.crt | awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' | base64 -w 0)
CERT_B64=$(cat ${CERTS_DIR}/${AGENT_NAME}.crt | awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' | base64 -w 0)
PRIV_KEY_B64=$(cat ${CERTS_DIR}/${AGENT_NAME}.key | awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' | base64 -w 0)
KEY_PASS=$(cat ${CERTS_DIR}/${AGENT_NAME}.pass)

LDFLAGS=(
    "-X main.AgentName=${AGENT_NAME}"
    "-X main.ServerAddr=${SERVER_ADDR}"
    "-X main.CACertB64=${CA_CERT_B64}"
    "-X main.CertB64=${CERT_B64}"
    "-X main.PrivKeyB64=${PRIV_KEY_B64}"
    "-X main.KeyPass=${KEY_PASS}"
    "-X main.IDHeader=${ID_HEADER}"
    "-X main.CmdHeader=${CMD_HEADER}"
    "-X main.DataHeader=${DATA_HEADER}"
    "-X main.RequestDelay=${REQUEST_DELAY}"
)
go build -o agents/${AGENT_NAME} -ldflags="${LDFLAGS[*]}" ./cmd/goc2agent
