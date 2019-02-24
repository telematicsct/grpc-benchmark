# gRPC Http mTLS Hmac benchmarking


## gRPC MTLS


### Run Server
`go run server/mtls/*.go --tls-cert=certs/server-cert.pem --tls-key=certs/server-key.pem`

`docker build -t dcm-server -f mtls.Dockerfile .`

`docker run -it dcm-server:latest --tls-cert=server-cert.pem --tls-key=server-key.pem`

### Run Client

`go run client/mtls/main.go --tls-cert=certs/server-cert.pem`
