FROM golang:1.11.1-alpine3.8 as builder
LABEL stage=intermediate

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go
ENV GO111MODULE=on

RUN set -x \
    && apk add --update --no-cache --virtual .build-deps \
    git \
    ca-certificates \
    gcc \
    libc-dev \
    libgcc \
    make \
    && apk add --no-cache upx

# To create a rootless container
RUN adduser -D -g '' gkuser

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY Makefile .
COPY dcm/* ./dcm/
COPY server/mtls/*.go ./

COPY ./certs/server-cert.pem /go/bin/
COPY ./certs/server-key.pem /go/bin/

RUN make clean \
    && make static \
    && apk del .build-deps \
    && echo "Build complete."

# Compress go binary
# https://linux.die.net/man/1/upx
RUN upx -7 -qq output/dcm-server && \
    upx -t output/dcm-server && \
    mv output/dcm-server /go/bin/dcm-server

#gcr.io/distroless/base
FROM scratch

WORKDIR /app

COPY --from=builder /go/bin/dcm-server ./
COPY --from=builder /go/bin/server-cert.pem ./
COPY --from=builder /go/bin/server-key.pem ./

COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /etc/passwd /etc/passwd

USER gkuser

ENTRYPOINT [ "./dcm-server" ]
EXPOSE 7900
EXPOSE 7901