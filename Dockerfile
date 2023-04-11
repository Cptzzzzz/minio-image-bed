FROM golang:1.19 AS builder

WORKDIR /work

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

RUN go build -o image-server server.go


FROM alpine

WORKDIR /app

COPY --from=builder /work/image-server /app/image-server

EXPOSE 8080

CMD /app/image-server $Endpoint $AccessKey $SecretKey $Bucket