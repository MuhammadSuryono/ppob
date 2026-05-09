FROM golang:1.21-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/main .
COPY --from=builder /build/config ./config
COPY --from=builder /build/internal ./internal
COPY --from=builder /build/cmd ./cmd
COPY --from=builder /build/internal/handlers ./internal/handlers
COPY --from=builder /build/internal/services ./internal/services

ENV SERVER_PORT=8080

EXPOSE 8080

CMD ["./main"]