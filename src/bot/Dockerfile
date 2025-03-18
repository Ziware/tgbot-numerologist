FROM golang:1.24-alpine as builder

WORKDIR /app

COPY . .
RUN ls
RUN go mod download -x

WORKDIR /app/bot
RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bot .

ENTRYPOINT ["./bot"]
