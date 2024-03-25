FROM golang:1.21 AS builder

ENV GO111MODULE=on\
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum /
RUN go mod download

COPY . .
RUN go build -o app .

FROM scratch

COPY --from=builder /build/app /
EXPOSE 8888
ENTRYPOINT [ "/app" ]

