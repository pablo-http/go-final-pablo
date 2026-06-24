# Этап сборки
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .

# Финальный образ
FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/scheduler .
COPY web/ web/

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db

CMD ["./scheduler"]
