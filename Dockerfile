FROM golang:1.19 AS builder

WORKDIR /work

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download

RUN go build -o server server.go


FROM alpine

WORKDIR /app

COPY --from=builder /work/server /app/server

EXPOSE 8080

CMD /app/server $Endpoint $AccessKey $SecretKey $Bucket