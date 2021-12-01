openssl pkcs12 -in Certificates.p12 -nocerts -out private-key.pem -nodes
openssl pkcs12 -in Certificates.p12 -nokeys -out certificate.crt
