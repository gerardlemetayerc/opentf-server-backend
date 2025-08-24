FROM golang:1.23 as builder

WORKDIR /app

RUN apt-get update && apt-get install -y build-essential

RUN git clone https://github.com/gerardlemetayerc/opentf-server-backend.git .

RUN go mod tidy
RUN CGO_ENABLED=1 go build -o opentf-server ./cmd/main.go

# Utilise une image Debian pour l'exécution
FROM debian:stable-slim
WORKDIR /app

COPY --from=builder /app/opentf-server .

EXPOSE 8080

CMD ["./opentf-server"]