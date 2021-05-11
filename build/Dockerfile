FROM golang:1.16.0-buster

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...

CMD ["go", "run", "./cmd/server"]
