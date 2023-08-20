FROM golang:1.19 as builder

RUN  apt update && apt install -y git ca-certificates

WORKDIR /app

COPY . .
#RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cmd/trade -ldflags="-w -s" cmd/trade/main.go

CMD [ "tail", "-f", "/dev/null" ]