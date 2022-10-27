#!/bin/bash
CERTS_DIR="./certs/"
SUBJECT=${1:?"usage: ${0} <subject>"}
KEY_PASS=$(openssl rand -base64 32)

# All newly created files should have chmod 600
umask 077

# Create a RootCA certificate and private key if does not exist
if [[ ! -f "${CERTS_DIR}/ca.crt" ]]; then
    openssl rand -base64 32 > ${CERTS_DIR}/ca.pass
    openssl req -x509 -new -sha512 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 \
        -nodes -subj /CN=ca -days 730 -passin file:${CERTS_DIR}/ca.pass \
        -out ${CERTS_DIR}/ca.crt -keyout ${CERTS_DIR}/ca.key
fi

# Generate a password file for the private key
echo -n "${KEY_PASS}" > ${CERTS_DIR}/${SUBJECT}.pass

# Generate a password-protected private key
openssl ecparam -genkey -name prime256v1 | \
    openssl ec -aes256 -passout file:${CERTS_DIR}/${SUBJECT}.pass -out ${CERTS_DIR}/${SUBJECT}.key

# Generate a config
cat > ${CERTS_DIR}/${SUBJECT}.conf << EOF
[ req ]
prompt             = no
days               = 730
default_bits       = 4096
distinguished_name = req_distinguished_name
req_extensions     = req_ext

[ req_distinguished_name ]
countryName        = PL
localityName       = Warsaw
organizationName   = NSA
commonName         = ${SUBJECT}

[ req_ext ]
subjectAltName     = @alt_names

[alt_names]
DNS.1              = ${SUBJECT}
EOF

# Generate a certificate signing request
openssl req -new -key ${CERTS_DIR}/${SUBJECT}.key -nodes -config ${CERTS_DIR}/${SUBJECT}.conf \
    -passin file:${CERTS_DIR}/${SUBJECT}.pass -out ${CERTS_DIR}/${SUBJECT}.csr

# Generate a certificate
openssl x509 -req -in ${CERTS_DIR}/${SUBJECT}.csr -nodes -extensions req_ext -extfile ${CERTS_DIR}/${SUBJECT}.conf \
    -CA ${CERTS_DIR}/ca.crt -CAkey ${CERTS_DIR}/ca.key -CAcreateserial \
    -sha512 -out ${CERTS_DIR}/${SUBJECT}.crt 

# Remove the CSR and config files
rm -f ${CERTS_DIR}/${SUBJECT}.csr && rm -f ${CERTS_DIR}/${SUBJECT}.conf