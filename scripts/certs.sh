!#/bin/sh

mkdir tmp
openssl genrsa -out tmp/ca.key 4096
openssl req -new -x509 -days 1826 -key tmp/ca.key -out tmp/ca.crt -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com"
openssl genrsa -out tmp/server.key 4096
openssl req -new -key tmp/server.key -out tmp/server.csr -subj "/C=US/ST=Denial/L=Springfield/O=Dis/CN=www.example.com"
openssl x509 -req -days 730 -in tmp/server.csr -CA tmp/ca.crt -CAkey tmp/ca.key -set_serial 01 -out tmp/server.crt
