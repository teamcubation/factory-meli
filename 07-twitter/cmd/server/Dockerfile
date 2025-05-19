FROM golang:1.24.1-alpine AS builder
WORKDIR /app
COPY ../../ ./
WORKDIR /app/07-twitter/cmd/server
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/07-twitter/cmd/server/main .
EXPOSE 8080
CMD ["./main"]