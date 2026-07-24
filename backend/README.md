# Бэкенд

Go-сервер для игры «Детектив» на стандартном `net/http`.

## Запуск

Скопируй `.env.example` в `.env` и пропиши переменные:

```bash
cp .env.example .env
go run ./cmd/server/
```

Переменные окружения:

| Переменная | По умолчанию | Описание |
|---|---|---|
| `PORT` | `8080` | Порт сервера |
| `OPENROUTER_API_KEY` | — | API-ключ OpenRouter |
| `OPENROUTER_MODEL` | `openai/gpt-4o-mini` | Модель для запросов |
| `MONGO_URI` | `mongodb://localhost:27017` | Подключение к MongoDB |
| `MONGO_DATABASE` | `detective_game` | Имя базы данных |

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
├── infrastructure/ — адаптеры (in-memory storage, mock LLM)
└── interfaces/     — HTTP-обработчики, роутер, middleware
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

Все эндпоинты требуют заголовок `X-User-ID`. Эндпоинты с игровой логикой также требуют `X-Session-ID` (кроме `POST /sessions` и `GET /sessions/history`).

| Метод | Путь | Описание |
|-------|------|---------|
| POST | `/api/v1/sessions` | Создать игровую сессию |
| GET | `/api/v1/sessions/history` | История сессий |
| GET | `/api/v1/sessions/current` | Текущая сессия |
| GET | `/api/v1/characters` | Список персонажей |
| GET | `/api/v1/characters/{charId}` | Детали персонажа |
| GET | `/api/v1/evidence` | Список улик |
| GET | `/api/v1/evidence/{evId}` | Детали улики |
| GET | `/api/v1/chronology` | Хронология событий |
| PATCH | `/api/v1/chronology/{chronId}/notes/{noteId}` | Обновить заметки |
| POST | `/api/v1/interrogations` | Начать допрос |
| POST | `/api/v1/interrogations/{interId}/messages` | Отправить сообщение |
| GET | `/api/v1/interrogations/{interId}/messages` | История допроса |
| PATCH | `/api/v1/interrogations/{interId}/complete` | Завершить допрос |
| POST | `/api/v1/actions/dna-analysis` | Анализ ДНК |
| POST | `/api/v1/actions/fingerprints` | Отпечатки пальцев |
| POST | `/api/v1/actions/alibi-check` | Проверка алиби |
| POST | `/api/v1/actions/camera-review` | Записи с камер |
| POST | `/api/v1/actions/call-history` | История звонков |
| POST | `/api/v1/actions/transactions` | Банковские операции |
| POST | `/api/v1/reports` | Отправить финальный отчёт |
| GET | `/api/v1/reports/{reportId}` | Просмотр отчёта |

## Текущее состояние

- **Storage:** In-memory (сбрасывается при перезапуске)
- **LLM:** Mock-сервис с предзаданным сценарием (убийство на вилле)
- **Персонажи:** 5 предзагруженных прототипов
