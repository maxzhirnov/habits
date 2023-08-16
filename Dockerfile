# Этап сборки
FROM golang:1.20 as builder

WORKDIR /app

# Копирование модулей и их установка
COPY go.mod go.sum ./
RUN go mod download

COPY start.sh /start.sh
RUN chmod +x /start.sh

COPY config/config.yml /config.yml


# Копирование исходного кода приложения
COPY . .

# Сборка cmd/client/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /client ./cmd/client/main.go

# Сборка cmd/habits/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /habits ./cmd/habits/main.go

# Этап выполнения
FROM alpine:latest

WORKDIR /

COPY --from=builder /client /client
COPY --from=builder /habits /habits
COPY --from=builder /start.sh /start.sh
COPY --from=builder /config.yml /config/config.yml

CMD ["./start.sh"]
