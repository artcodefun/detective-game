# ДетектИИв

Мобильная детективная игра: игрок изучает дело, допрашивает подозреваемых,
проверяет улики и сдаёт итоговый отчёт. Сервер управляет расследованиями,
хранит игровое состояние и генерирует сценарии и ответы персонажей через LLM.

## Состав проекта

- `mobile/` — Flutter-клиент для Android и iOS.
- `backend/` — Go API, MongoDB-хранилище и интеграция с OpenRouter.
- `backend/api/openapi-v1.yaml` — актуальный контракт HTTP API.

## Требования

- FVM;
- Flutter 3.44.1, установленный через FVM;
- Go 1.25 или новее;
- доступная MongoDB;
- ключ OpenRouter для генерации игрового контента.

## Быстрый старт

```bash
cp backend/.env.example backend/.env
cd backend
set -a
source .env
set +a
go run ./cmd/server
```

В другом терминале запустите приложение, передав адрес API:

```bash
cd mobile
fvm flutter pub get
fvm flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

Для Android Emulator обычно нужен адрес `http://10.0.2.2:8080`, для iOS
Simulator — `http://localhost:8080`.

Проверки перед коммитом:

```bash
(cd backend && go test ./...)
(cd mobile && fvm flutter analyze && fvm flutter test)
```

Подробности находятся в [`backend/README.md`](backend/README.md) и
[`mobile/README.md`](mobile/README.md).

## Лицензия

Проект распространяется по лицензии MIT. См. [`LICENSE`](LICENSE).
