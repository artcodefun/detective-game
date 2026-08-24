# Бэкенд

Go-сервер для игры «Детектив» на стандартном `net/http`.

## Запуск

Скопируй `.env.example` в `.env`, пропиши переменные и экспортируй их в
окружение текущей оболочки перед запуском:

```bash
cp .env.example .env
set -a
source .env
set +a
go run ./cmd/server
```

## Docker

Compose запускает API и MongoDB, а данные базы сохраняет в именованном томе
`mongo_data`. Из каталога `backend/`:

```bash
cp .env.example .env
# Заполни OPENROUTER_API_KEY и задай надёжный MONGO_ROOT_PASSWORD.
docker compose up -d mongo
set -a
source .env
set +a
docker compose exec mongo mongosh --quiet --username "$MONGO_ROOT_USERNAME" --password "$MONGO_ROOT_PASSWORD" --authenticationDatabase admin --eval 'rs.initiate({_id: "rs0", members: [{_id: 0, host: "mongo:27017"}]})'
docker compose up -d --build api
```

`rs.initiate()` выполняется один раз. Проверка:

```bash
docker compose exec mongo mongosh --quiet --username "$MONGO_ROOT_USERNAME" --password "$MONGO_ROOT_PASSWORD" --authenticationDatabase admin --eval 'rs.status().set'
```

API будет доступен только с сервера по `http://127.0.0.1:${HOST_PORT:-8080}`.
Для внешнего доступа используй reverse proxy.

### MongoDB Compass

MongoDB публикуется только на loopback-интерфейсе сервера. Для доступа к ней с
локального компьютера подними SSH-туннель:

```bash
ssh -N -L 27017:127.0.0.1:27017 <user>@<server>
```

После этого подключись в MongoDB Compass по строке:

```text
mongodb://detective_game:<MONGO_ROOT_PASSWORD>@127.0.0.1:27017/?authSource=admin&directConnection=true
```

Если `27017` занят только на локальном компьютере, измени левую часть
SSH-туннеля, например `-L 27018:127.0.0.1:27017`. Если меняешь
`MONGO_HOST_PORT` на сервере, используй это значение в правой части туннеля.

Полезные команды:

```bash
docker compose logs -f api
docker compose down
```

Переменные окружения:

| Переменная | По умолчанию | Описание |
|---|---|---|
| `PORT` | `8080` | Порт сервера |
| `OPENROUTER_API_KEY` | — | API-ключ OpenRouter |
| `OPENROUTER_MODEL` | `deepseek/deepseek-v4-flash` | Модель для запросов |
| `MONGO_URI` | `mongodb://localhost:27017/?replicaSet=rs0` | Подключение к MongoDB |
| `MONGO_DATABASE` | `detective_game` | Имя базы данных |
| `IOS_MIN_SUPPORTED_VERSION` | `0.0.0` | Минимальная версия приложения для iOS |
| `ANDROID_MIN_SUPPORTED_VERSION` | `0.0.0` | Минимальная версия приложения для Android |
| `IOS_UPDATE_URL` | — | Ссылка на обновление приложения для iOS |
| `ANDROID_UPDATE_URL` | — | Ссылка на обновление приложения для Android |

## Архитектура

Сервис использует **Hexagonal Architecture (Ports and Adapters)**, **Domain-Driven Design (DDD)**
и **CQRS** (Command Query Responsibility Segregation).

```
cmd/server/         — точка входа
internal/
├── domain/         — агрегаты, value objects, доменные enum-ы
├── application/    — слой приложения
│   ├── commands.go / queries.go  — публичные CQRS-интерфейсы
│   ├── errors.go                 — AppError, ErrorKind
│   ├── actor.go                  — контекст запроса (Actor)
│   ├── readmodels/               — read-модели + мапперы domain → readmodel
│   ├── commands/                 — реализации команд
│   ├── queries/                  — реализации запросов
│   ├── ports/                    — интерфейсы инфраструктуры
│   └── services/                 — сервисы приложения
├── bootstrap/      — композиция зависимостей (Pure DI)
├── infrastructure/ — адаптеры (MongoDB, OpenRouter LLM)
├── interfaces/     — HTTP-обработчики, роутер, middleware
api/                — OpenAPI-спецификация (openapi-v1.yaml)
```

## Правила

- **Агрегаты инкапсулированы.** Нельзя обращаться к полям агрегатов напрямую
  извне `domain/`. Любое изменение состояния — только через доменные методы агрегата.
- **Команды не возвращают доменные структуры.** Результат команды — скалярное значение
  или ID. Данные для чтения отдаются через запросы.
- **Запросы возвращают только read-модели.** Запросы никогда не возвращают доменные
  структуры. Все данные для внешнего потребителя проходят через `readmodels/`.
- **Порты определяются внутри.** Интерфейсы репозиториев и внешних сервисов живут
  в `application/ports/` и ссылаются только на доменные типы.
- **Реализации — снаружи.** Инфраструктурные адаптеры (`infrastructure/`) реализуют
  порты, но домен и приложение о них не знают.

## API

При первом запуске приложение вызывает `POST /v1/auth/anonymous` и получает Bearer token. Все защищённые эндпоинты требуют заголовок `Authorization: Bearer <token>`. Большинство эндпоинтов игровой логики также требуют `X-Session-ID` (UUID сессии). Исключения (не требуют `X-Session-ID`): `POST /sessions`, `GET /sessions/history`, `GET /sessions/current`, `GET /sessions/{id}`.

| Метод | Путь | Описание |
|-------|------|---------|
| POST | `/v1/sessions` | Создать игровую сессию |
| GET | `/v1/sessions/current` | Активная сессия пользователя (не требует `X-Session-ID`) |
| GET | `/v1/sessions/{id}` | Конкретная сессия по ID (не требует `X-Session-ID`) |
| GET | `/v1/sessions/history` | История завершённых сессий |
| GET | `/v1/characters` | Список персонажей |
| GET | `/v1/characters/{charId}` | Детали персонажа |
| GET | `/v1/evidence` | Список улик |
| GET | `/v1/evidence/{evId}` | Детали улики |
| GET | `/v1/chronology` | Хронология событий |
| PATCH | `/v1/chronology/{chronId}/notes/{noteId}` | Обновить заметки |
| POST | `/v1/interrogations` | Начать допрос |
| POST | `/v1/interrogations/{interId}/messages` | Отправить сообщение |
| GET | `/v1/interrogations/{interId}/messages` | История допроса |
| PATCH | `/v1/interrogations/{interId}/complete` | Завершить допрос |
| POST | `/v1/actions/dna-analysis` | Анализ ДНК |
| POST | `/v1/actions/fingerprints` | Отпечатки пальцев |
| POST | `/v1/actions/alibi-check` | Проверка алиби |
| POST | `/v1/actions/camera-review` | Записи с камер |
| POST | `/v1/actions/call-history` | История звонков |
| POST | `/v1/actions/transactions` | Банковские операции |
| POST | `/v1/reports` | Отправить финальный отчёт |
| GET | `/v1/reports/{reportId}` | Просмотр отчёта |

Полная спецификация в формате OpenAPI 3.1: [`api/openapi-v1.yaml`](api/openapi-v1.yaml).

## Проверки

```bash
gofmt -w .
go test ./...
```
