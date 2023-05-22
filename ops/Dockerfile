FROM golang:1.20.4-alpine3.18 AS builder
WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ttto ./cmd/ttto

FROM scratch as prod

COPY --from=builder /app/ttto /bin/ttto
