FROM golang:1.25.3 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o grpc_server ./storage

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/grpc_server .

EXPOSE 50000
ENTRYPOINT ["./grpc_server"]