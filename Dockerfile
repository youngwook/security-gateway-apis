FROM golang:1.12-alpine AS builder

ENV GO111MODULE=on

WORKDIR /go/src/github.com/youngwook/security-gateway-apis

RUN mkdir -pv vault/config

RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories

RUN apk update && apk add make git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go get && CGO_ENABLED=0 GOOS=linux go build 

FROM scratch

WORKDIR /

EXPOSE 8088

COPY --from=builder /go/src/github.com/youngwook/security-gateway-apis/res/configuration.toml /res/configuration.toml

COPY --from=builder /go/src/github.com/youngwook/security-gateway-apis/vault/config /vault/config

COPY --from=builder  /go/src/github.com/youngwook/security-gateway-apis .

ENTRYPOINT ["./security-gateway-apis"]
