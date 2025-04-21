# 🛡️ MEDODS Auth Service

Сервис аутентификации на Go.  
Реализует выдачу и обновление Access/Refresh токенов, хранение хеша refresh в PostgreSQL, защиту от повторного использования и контроль IP-адреса.

---

## ⚙️ Технологии

- Go 1.23.2
- PostgreSQL 16
- JWT (алгоритм SHA512)
- bcrypt + SHA256 (refresh token)
- Chi Router
- Docker + Docker Compose
- Unit-тесты (`go test`)

---

## 📌 Возможности

✅ Выдаёт пару `Access` / `Refresh` токенов по GUID  
✅ Хэширует refresh токен (SHA256 → bcrypt) и сохраняет в PostgreSQL  
✅ Один refresh можно использовать только один раз  
✅ При изменении IP логируется предупреждение (моковая отправка email)  
✅ Access и Refresh токены обоюдно связаны  
✅ Access хранится на клиенте, Refresh — в `HttpOnly` куки  
✅ Поддержка тестов логики обновления и защиты от повторного использования

---

## 📤 REST API

| Метод | URL             | Описание                                  |
|-------|------------------|-------------------------------------------|
| POST  | `/api/create`    | Выдаёт пару токенов (GUID в теле запроса) |
| PUT   | `/api/refresh`   | Обновляет токены по refreshToken (из куки) |

---

## 🔐 Форматы токенов

### Access token (JWT):
- Алгоритм: `HS512`
- Payload: `guid`, `ip`, `exp`
- Хранение: только у клиента (`Authorization: Bearer ...`)

### Refresh token:
- Формат: JWT → SHA256 → bcrypt
- Хранение: в базе данных только хэш
- Передаётся: как `HttpOnly` кука (base64)

---

## 🧪 Пример запроса

### POST `/api/create`

#### Тело запроса:
```json
{
  "guid": "123e4567-e89b-12d3-a456-426614174000"
}
```

#### Ответ:
- Заголовок: `Authorization: Bearer <access_token>`
- Кука: `refreshToken=<base64>; HttpOnly; Path=/api/refresh`

---

## 🐳 Docker + Docker Compose

### Dockerfile

```dockerfile
FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
RUN go build -o server cmd/server/main.go

FROM scratch
COPY --from=builder /app/server /server
EXPOSE 8901
ENTRYPOINT ["/server"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  db:
    image: postgres:16
    container_name: postgres_MEDODS
    restart: always
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: MEDODS
      POSTGRES_USER: senya
      POSTGRES_PASSWORD: 1q2w3e4r5t
    healthcheck:
        test: [ "CMD-SHELL", "pg_isready -U senya -d postgres" ]
        interval: 5s
        retries: 5

  api:
    container_name: api_MEDODS
    build: .
    depends_on:
      db:
        condition: service_healthy
    environment:
      ADDRESS: ":8901"
      DATABASE_HOST: "db"
      DATABASE_PORT: "5432"
      DATABASE_USER: "senya"
      DATABASE_PASSWORD: "1q2w3e4r5t"
      DATABASE_NAME: "MEDODS"
      DATABASE_SSL_MODE: "disable"
      ACCESS_KEY: "f4pq3792h3dy4g82o63R84P265o3874wgfiy2p947gf7qo5hcnbvtbo8y2c9upnox3q9E3"
      REFRESH_KEY: "x53416ucertiyvuybiunb5yp6no78iu65cr34exqyto839p28u320kjfnubviry3294bdf"
      LIFETIME_ACCESS: "15m"
      LIFETIME_REFRESH: "24h"
    ports:
      - "8901:8901"
```

---

## ✅ Тестирование

```bash
go test ./internal/server/app/auth
```

---

## 🧱 Структура проекта

```
cmd/                 # main.go
internal/
├── config/          # конфигурация (env, mock)
├── server/
│   ├── app/         # бизнес-логика
│   ├── delivery/    # HTTP handlers и middleware
│   ├── repository/  # хранилище PostgreSQL
pkg/
├── jwttoken/        # работа с JWT
```

---

## 📝 Автор

Senya 🚀  
Тестовое задание для MEDODS (Backend Developer)

---