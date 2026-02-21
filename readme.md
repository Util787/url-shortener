## Краткая документация API

Небольшая инструкция по запуску и использованию HTTP API сервиса сокращения ссылок.

### Коротко
- Базовый префикс: `/api/v1`
- Формат: JSON для POST-запросов

### Переменные окружения
Скопируйте `.env.example` в `.env` и заполните значения перед запуском.

Основные переменные (в файле `.env.example`):

- `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` — доступ к PostgreSQL
- `HTTP_SERVER_HOST`, `HTTP_SERVER_PORT` — адрес и порт HTTP сервера
- `TG_BOT_TOKEN` — токен телеграм-бота (опционально)

### Запуск
1. Запустите PostgreSQL (например, в Docker) и примените миграции из `migrations/postgres`.
2. Создайте `.env` и заполните переменные.
3. Запустите сервис:

```bash
go run ./cmd
```

(или `go build` и запуск сгенерированного бинарника)

### Эндпоинты
1) Сохранить длинную ссылку

- URL: `POST /api/v1/save`
- Body (JSON):

```json
{ "long_url": "https://example.com/some/very/long/path" }
```
- Успешный ответ: 200

```json
{ "message": "URL saved successfully" }
```

2) Получить случайную короткую ссылку

- URL: `GET /api/v1/random`
- Успешный ответ: 200

```json
{ "short_url": "abc123" }
```

3) Перенаправление по короткой ссылке

- URL (intended): `GET /api/v1/redirect/:short_url`
- Описание: возвращает редирект (HTTP 302) на соответствующую длинную ссылку.

Пример (curl):

```bash
curl -v http://localhost:8080/api/v1/redirect/abc123
```

4) Удалить короткую ссылку

- URL (intended): `POST /api/v1/delete/:short_url`
- Пример (curl):

```bash
curl -X POST http://localhost:8080/api/v1/delete/abc123
```

### Замечание о несовпадении маршрутов
В коде обработчики читают параметр `short_url` через `c.Param("short_url")`, тогда как маршруты в `internal/adapters/rest/routes.go` зарегистрированы как `/redirect` и `/delete` без явного параметра в пути. Рекомендуется привести маршруты и обработчики в соответствие. Два простых варианта:

1. Изменить маршруты на `/redirect/:short_url` и `/delete/:short_url` (рекомендуется).
2. Или изменить обработчики на чтение из query-параметра или тела запроса, если так задуман API.

### Примеры curl

Save:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"long_url":"https://example.com"}' http://localhost:8080/api/v1/save
```

Get random:

```bash
curl http://localhost:8080/api/v1/random
```

Redirect (browser / curl):

```bash
curl -v http://localhost:8080/api/v1/redirect/abc123
```

Delete:

```bash
curl -X POST http://localhost:8080/api/v1/delete/abc123
```

---


