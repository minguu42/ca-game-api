FROM golang:1.15.8-buster

WORKDIR /go/src/ca-game-api
COPY . .

RUN go get -u github.com/go-sql-driver/mysql

CMD ["go", "run", "main.go"]
