FROM golang:1.19 AS builder

WORKDIR /work

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

RUN go run server.go $Endpoint $AccessKey $SecretKey $Bucket
#FROM alpine
#
#WORKDIR /app
#
#COPY --from=builder /work/server /app/server
#
#EXPOSE 8080

