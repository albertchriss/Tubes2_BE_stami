FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

ENV GOPROXY https://goproxy.io,direct

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/*.go

FROM alpine:3.21

RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /app/api .

CMD ["./api"]