FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates \
	&& adduser -D -u 10001 appuser

COPY --from=builder /out/api /app/api

ENV SERVER_HOST=0.0.0.0 \
	SERVER_PORT=8000 \
	FRONTEND_URL=http://localhost:3000 \
	COOKIE_SECURE=false \
	SESSION_DURATION_MINUTES=1440

EXPOSE 8000

USER appuser

CMD ["/app/api"]