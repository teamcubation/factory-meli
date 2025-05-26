FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test -v ./...

RUN go build -o app ./cmd

ENTRYPOINT ["./app"]