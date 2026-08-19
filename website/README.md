# Сайт «ДетектИИв»

Статический лендинг и публичные документы игры. После сборки содержимое
`out/` можно целиком разместить за nginx без Node.js на production-сервере.

## Локальный запуск

```bash
npm install
npm run dev
```

Сайт откроется на `http://localhost:3000`.

## Production-сборка

```bash
npm ci
npm run build
```

Готовые файлы появятся в `out/`. Пример публикации:

```bash
rsync -az --delete out/ user@server:/var/www/detective-game/
```

Для nginx используется `root /var/www/detective-game;` и
`try_files $uri $uri/ =404;`.

## Docker

Production-контейнер содержит только nginx и готовые статические файлы:

```bash
docker compose up -d --build
```

По умолчанию сайт будет доступен на сервере по адресу
`http://127.0.0.1:8081`. Внешний nginx должен перенаправлять на него запросы
домена `detective-game.artcodefun.com`. Порт можно изменить переменной
`WEBSITE_PORT`.

Перед первым запуском на сервере создайте локальный файл окружения:

```bash
cp .env.example .env
```

Например, для публикации сайта на `127.0.0.1:7071`:

```dotenv
WEBSITE_PORT=7071
```

Docker Compose автоматически прочитает `.env`, если команда запускается из
каталога `website/`. Сам файл `.env` не добавляется в Git.

## Настройки

Название разработчика, email поддержки, адрес сайта и дата редакции документов
находятся в `app/site-config.ts`.
