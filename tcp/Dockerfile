FROM golang:1-alpine as builder

RUN apk --no-cache --no-progress add git ca-certificates tzdata make \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /go/whoami

COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on GOPROXY=https://proxy.golang.org go mod download

COPY . .

RUN make build

FROM golang:1-alpine
RUN apk add --no-cache netcat-openbsd

COPY --from=builder /go/whoami/whoami /whoami

ENTRYPOINT ["/whoami"]
EXPOSE 8080
