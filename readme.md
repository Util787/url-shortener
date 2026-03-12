# 🔗 URL Shortener Service

Микросервис для сокращения длинных ссылок.  
Сервис сохраняет ссылки в PostgreSQL, отдаёт их по REST API и поддерживает редиректы.

## Features
- ✅ REST API для сохранения и получения ссылок
- ✅ Редирект по короткой ссылке
- ✅ Поддержка PostgreSQL
- ✅ Telegram-бот (опционально)
- ✅ Готов к запуску через Docker Compose

---

## Quick start 🚀

### Requirements 📦
- Docker
- Go 1.24+ (только для ручного запуска)

---

### 1. Clone Repository 📂
```bash
git clone https://github.com/Util787/url-shortener && cd url-shortener
```

---

### 2. Configure .env ⚙️
Скопируйте `.env.example` в `.env` и настройте переменные окружения:

```bash
cp .env.example .env
```

---

### 3. Run with Docker Compose 🐳
```bash
docker compose up --build
```
Миграции применятся автоматически.

---

## Manual run (without Docker) 🛠️
Если хотите запустить локально:

```bash
go mod tidy
go run ./cmd/main.go
```

Не забудьте предварительно запустить PostgreSQL и применить миграции из `migrations/postgres`.

---

## 🌐 API Endpoints

### `POST /save`
Сохраняет длинную ссылку и возвращает короткую.

**Request:**
```json
{
  "long_url": "https://example.com/very/long/url"
}
```

**Response 200 OK:**
```json
{
  "message": "URL saved successfully with short URL: http://localhost:8080/abc123"
}
```

---

### `GET /random`
Возвращает случайную короткую ссылку.

**Response 200 OK:**
```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

---

### `GET /:short_url_id`
Редирект на оригинальную длинную ссылку.

**Example:**
```
GET /abc123
```
- ✅ Если найдена — HTTP 302 Redirect
- ❌ Если нет — 404 Not Found

---

### `POST /delete`
Удаляет короткую ссылку.

**Request:**
```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

**Response 200 OK:**
```json
{
  "message": "URL deleted successfully"
}
```
