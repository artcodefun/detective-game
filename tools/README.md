# Инструменты генерации контента

Jupyter-ноутбуки и промпты для генерации игровых ассетов. Запускаются на Yandex Cloud с GPU.

## Структура

```
tools/
├── characters/    # Портреты персонажей (Stable Diffusion)
├── evidence/      # Изображения улик (Stable Diffusion)
├── audio/         # Аудио-ассеты с псевдоязыком
└── README.md
```

## Рабочий процесс

1. Заполнить `characters-sd-gen.json` промптами
2. Открыть `characters.ipynb` в Yandex Cloud JupyterLab
3. Запустить все ячейки — на выходе папка с PNG
4. Скопировать результат в `mobile/assets/characters/`

## Формат `*-sd-gen.json`

```json
[
  {
    "id": "char_01",
    "name": "Иван Петров",
    "filename": "ivan_petrov.png",
    "item_prompt": "realistic portrait, middle-aged man, detective story style, neutral expression"
  }
]
```
