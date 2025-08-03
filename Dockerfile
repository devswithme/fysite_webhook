FROM golang:tip-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o webhook .

FROM alpine:latest

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/webhook .

RUN chown -R appuser /app

USER appuser

EXPOSE 8080

CMD ["./webhook"]