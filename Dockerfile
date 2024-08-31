FROM golang:1.22.2 as builder

WORKDIR /app

# Install required packages
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o echo-server

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/echo-server /echo-server
COPY --from=builder /app/assets /assets

EXPOSE 8080

ENTRYPOINT ["/echo-server"]