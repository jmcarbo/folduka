openssl pkcs12 -in $1 -out temp.pem -passin pass:$2 -nodes
openssl pkcs12 -export -in temp.pem -out $1 -passout pass:$3
rm -rf temp.pem
