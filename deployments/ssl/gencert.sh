#!/bin/bsah

initca() {
  cfssl gencert -initca ./ca-csr.json | cfssljson -bare ca 
}

gencert() {
  cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=./ca-config.json \
        -profile=server ./server-csr.json | cfssljson -bare 
}