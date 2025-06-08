FROM golang:onbuild

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/url-shortener .
COPY --from=builder /app/web ./web
COPY --from=builder /app/internal/config/local.yaml ./config.yaml

EXPOSE 8080

CMD ["./url-shortener"]