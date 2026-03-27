FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sliver-web ./cmd/server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/sliver-web .
COPY --from=builder /app/sliver-web.conf.example .

RUN mkdir -p /data /certs

EXPOSE 8080

ENV SLIVER_GRPC_ADDR=sliver-server:50050
ENV SLIVER_CERTS_DIR=/certs
ENV DB_PATH=/data/sliver-web.db
ENV LOG_LEVEL=info

ENTRYPOINT ["./sliver-web"]
CMD ["--config", "sliver-web.conf"]
