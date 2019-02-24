PROJECT_ROOT=github.com/telematicsct/grpc-benchmark

# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)

NAME 		:= dcm-service
PKG 		:= github.com/telematicsct/grpc-benchmark
BUILDTAGS 	:=
DISTDIR 	:= ${PREFIX}/output

.PHONY: clean
clean: ## Cleanup any build binaries or packages
	$(RM) $(NAME)
	$(RM) -r $(DISTDIR)
	

.PHONY: static
static:
	GOOS=linux CGO_ENABLED=0 go build -tags "$(BUILDTAGS) static_build" -ldflags "-linkmode internal -extldflags -static" \
    -o output/dcm-server *.go

.PHONY: proto
proto:
	protoc -I dcm/ dcm/dcm.proto --go_out=plugins=grpc:dcm

.PHONY: gencert
gencert:
	rm -rf certs/*.pem
	openssl req -config certs/mtls.conf -new -x509 -newkey rsa:2048 -nodes -keyout certs/server-key.pem -days 3650 -out certs/server-cert.pem
	openssl x509 -in certs/server-cert.pem -text -noout

.PHONY: image
image: ## Creates the docker images of the app and cleanups the intermediate
	echo '>> build docker image'
	docker build -t telematicsct/dcm-server .
	docker image prune --force --filter label=stage=intermediate