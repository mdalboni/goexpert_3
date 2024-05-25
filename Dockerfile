FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
RUN CGO_ENABLED=0 GOOS=linux go build -o cloudrun ./cmd/main.go

FROM scratch
WORKDIR /app
COPY --from=builder /app/cloudrun .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT [ "./cloudrun" ]