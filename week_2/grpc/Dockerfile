FROM golang:1.21 AS builder

COPY . /github.com/FreylGit/microservices/week_2/grpc/sourse/
WORKDIR  /github.com/FreylGit/microservices/week_2/grpc/sourse/

RUN go mod download
RUN go build -o ./bin/crud_server ./cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/FreylGit/microservices/week_2/grpc/sourse/bin/crud_server .


CMD ["./crud_server"]
CMD ["ls"]