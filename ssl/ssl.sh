#!/bin/bash
# Inspired by: https://github.com/grpc/grpc-java/tree/master/examples#generating-self-signed-certificates-for-use-with-grpc
# and https://github.com/grpc/grpc-java/tree/master/examples/example-tls

# Output files will be stored in the ./ssl directory.
SSL_DIR="./ssl"
mkdir -p "$SSL_DIR"

# Define file paths
CA_KEY="$SSL_DIR/ca.key"
CA_CRT="$SSL_DIR/ca.crt"
SERVER_KEY="$SSL_DIR/server.key"
SERVER_CSR="$SSL_DIR/server.csr"
SERVER_CRT="$SSL_DIR/server.crt"
SERVER_PEM="$SSL_DIR/server.pem"

# Change this CN to match your host in your environment if needed.
SERVER_CN="localhost"

# Create a configuration file for the extensions
cat <<EOF > "$SSL_DIR/server_cert_ext.cnf"
[req]
distinguished_name = req_distinguished_name
req_extensions = req_ext
x509_extensions = v3_req
[req_distinguished_name]
countryName = Country Name (2 letter code)
countryName_default = US
stateOrProvinceName = State or Province Name (full name)
stateOrProvinceName_default = California
localityName = Locality Name (eg, city)
localityName_default = San Francisco
organizationName = Organization Name (eg, company)
organizationName_default = My Company
commonName = Common Name (e.g. server FQDN or YOUR name)
commonName_default = ${SERVER_CN}
[req_ext]
subjectAltName = @alt_names
[v3_req]
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${SERVER_CN}
EOF

# Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out "$CA_KEY" 4096
openssl req -passin pass:1111 -new -x509 -days 3650 -key "$CA_KEY" -out "$CA_CRT" -subj "/CN=${SERVER_CN}"

# Step 2: Generate the Server Private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out "$SERVER_KEY" 4096

# Step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key "$SERVER_KEY" -out "$SERVER_CSR" -subj "/CN=${SERVER_CN}" -config "$SSL_DIR/server_cert_ext.cnf"

# Step 4: Sign the certificate with the CA we created (it's called self-signing) - server.crt
openssl x509 -req -passin pass:1111 -days 3650 -in "$SERVER_CSR" -CA "$CA_CRT" -CAkey "$CA_KEY" -set_serial 01 -out "$SERVER_CRT" -extensions v3_req -extfile "$SSL_DIR/server_cert_ext.cnf"

# Step 5: Convert the server certificate to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in "$SERVER_KEY" -out "$SERVER_PEM"

echo "Certificates created successfully in $SSL_DIR."
