FROM golang:1.11.4-alpine3.8 as builder
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

COPY . .

RUN make \
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
COPY --from=builder /app/certs/* ./certs/

COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /etc/passwd /etc/passwd

USER gkuser

ENTRYPOINT [ "./dcm-server", "all", "--key=certs/server.crt", "--key=certs/server.key", "--ca=certs/ca.crt" ]
EXPOSE 7900
EXPOSE 7901