FROM golang:1.16.0-buster

WORKDIR /go/src/github.com/minguu42/ca-game-api
COPY . .

RUN go mod tidy

CMD ["go", "run", "cmd/main/main.go"]
