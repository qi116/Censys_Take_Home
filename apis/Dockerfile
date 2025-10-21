FROM golang:1.25.3 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rest_server ./apis

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/rest_server .

EXPOSE 8080
ENTRYPOINT ["./rest_server"]