# GeoLogic API

Backend-сервис для хранения пользователей и объектов аренды с геокоординатами, а также для поиска доступных объектов рядом с заданной точкой.

## Кратко о проекте

GeoLogic API предоставляет REST API для:
- управления пользователями;
- управления объектами аренды с координатами и описанием;
- поиска объектов аренды в заданном радиусе вокруг точки с помощью PostGIS.

Архитектура приложения строится по слоям: handler → service → repository с DI и централизованной обработкой ошибок.

## Стек технологий

- Go 1.26.1
- Gin
- GORM
- PostgreSQL + PostGIS
- log/slog
- github.com/caarlos0/env + godotenv
- testcontainers-go для e2e
- stretchr/testify для тестов

## Схема данных

### Users
- id — PK
- tg_id — bigint, unique, not null

### Rents
- id — PK
- location — geometry(Point, 4326)
- address — varchar
- info — varchar, nullable

## Что уже реализовано

- Базовая структура приложения с разделением на internal/api, internal/handler, internal/service, internal/repository и internal/database.
- CRUD для пользователей.
- CRUD для объектов аренды.
- Геопространственный поиск объектов аренды рядом с точкой через PostGIS.
- Централизованная обработка ошибок и middleware для HTTP-ответов.
- Логирование через slog.
- E2E-тесты для пользователей, аренд и геопоиска.
- OpenAPI-спецификация для текущего API.

## Текущий API

Базовый путь приложения:
- /api/users
- /api/rents
- /api/rents/available

## Что планируется дальше

- расширить модель домена и добавить новые сущности;
- покрыть сервисы и репозитории unit-тестами;
- улучшить OpenAPI и примеры запросов;
- добавить авторизацию и дополнительные бизнес-правила.

## Конфигурация и запуск

1. Скопировать шаблон окружения:
```bash
cp .env.template .env
```

2. Настроить переменные окружения для базы данных, порта сервера и логирования.

3. Запустить PostgreSQL с PostGIS:
```bash
task db:up
```

4. Запустить приложение:
```bash
task run
```

5. Запустить e2e-тесты:
```bash
go test -tags=e2e ./internal/tests/e2e
```