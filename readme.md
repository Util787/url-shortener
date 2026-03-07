## Краткая документация API

Небольшая инструкция по запуску и использованию HTTP API сервиса сокращения ссылок.

### Коротко
- Базовый префикс: `/`
- Формат: JSON для POST-запросов

### Переменные окружения
Скопируйте `.env.example` в `.env` и заполните значения перед запуском.

Основные переменные (в файле `.env.example`):

- `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` — доступ к PostgreSQL
- `HTTP_SERVER_HOST`, `HTTP_SERVER_PORT` — адрес и порт HTTP сервера
- `TG_BOT_TOKEN` — токен телеграм-бота (опционально)
- `REDIRECT_BASE_URL` - url по которому происходит редирект

## Запуск


Вариант 1 — Docker Compose (миграции автоматические)

- Скопируйте файл `.env.example` в `.env` и заполните значения переменных окружения 

```bash
docker compose up --build
```

Вариант 2 — Локально (ручной запуск и применение миграций)

- Запустите PostgreSQL (например, в Docker) или используйте существующую БД.
- Примените миграции из каталога `migrations/postgres` к вашей базе данных.
- Скопируйте файл `.env.example` в `.env` и заполните значения переменных окружения 
- Запустите сервер вручную:

```bash
go run ./cmd/main.go
```


# Эндпоинты

Короткое описание HTTP‑эндпоинтов сервиса (маршруты находятся в корне сервера).

- POST /save
	- Body (JSON): { "long_url": "https://example.com/..." }
	- Ответ 200: { "message": "URL saved successfully with short URL: http://<host>:<port>/<id>" }

- GET /random
	- Ответ 200: { "short_url": "http://<host>:<port>/<id>" }

- GET /:short_url_id
	- Пример: GET /5vrTql7AOt
	- Возвращает HTTP 302 редирект на оригинальную длинную ссылку

- POST /delete
	- Body (JSON): { "short_url": "http://<host>:<port>/<id>" } (или используйте `/delete/:short_url` если предпочитаете path param)
	- Ответ 200: { "message": "URL deleted successfully" }

---



