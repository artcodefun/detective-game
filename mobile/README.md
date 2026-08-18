# Мобильное приложение «ДетектИИв»

Flutter-клиент игры для Android и iOS. Клиент отвечает за интерфейс, локальные
аудионастройки, голосовой ввод и отображение расследования; игровые данные
получает из backend API.

## Запуск

```bash
fvm flutter pub get
fvm flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

`API_BASE_URL` по умолчанию равен `http://localhost:8080`. При запуске на
Android Emulator используйте `http://10.0.2.2:8080`; физическому устройству
нужен доступный в локальной сети адрес компьютера.

## Структура

```text
lib/
├── models/       # модели ответов API
├── screens/      # экраны и локальное UI-состояние
├── services/     # HTTP-клиент, сессия и аудио
└── widgets/      # переиспользуемые виджеты
```

Глобальные зависимости и состояние сессии предоставляются через `provider`.
Краткоживущее состояние конкретного экрана остаётся внутри его `State`.

## Проверки

```bash
fvm dart format --output=none --set-exit-if-changed lib test
fvm flutter analyze
fvm flutter test
```
