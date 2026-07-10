# Detective Game

Мобильная игра в жанре детективного расследования. Игрок допрашивает подозреваемых, анализирует улики и раскрывает преступления. Сессии генерируются LLM (DeepSeek).

## Структура

```
detective-game/
├── mobile/       # Flutter-приложение
├── backend/      # Бэкенд (на этапе продакшена)
├── tools/        # Инструменты генерации контента
└── SPEC.md       # Спецификация проекта
```

## Быстрый старт

```bash
# Мобильное приложение
cd mobile
flutter pub get
flutter run

# Генерация контента (на Yandex Cloud)
# Открыть tools/characters/characters.ipynb в Jupyter
```
