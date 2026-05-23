FROM golang:1.25-alpine AS builder

WORKDIR /build/app

RUN apk add --no-cache git ca-certificates

COPY wallet-service/go.mod wallet-service/go.sum ./
COPY shared ../shared
RUN go mod download

COPY wallet-service/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/app/main .
COPY --from=builder /build/app/config ./config
COPY --from=builder /build/app/internal ./internal
COPY --from=builder /build/app/cmd ./cmd

ENV SERVER_PORT=8080

EXPOSE 8080

CMD ["./main"]