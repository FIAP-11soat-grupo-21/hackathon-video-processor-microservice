FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

FROM alpine:latest

RUN apk --no-cache add ca-certificates ffmpeg

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/worker .

CMD ["./api"]
