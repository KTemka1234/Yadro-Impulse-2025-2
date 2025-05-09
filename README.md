# YADRO-IMPULSE-2025-2

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://golang.org/)

Второе тестовое задание для стажировки Yadro Impulse 2025. Исходное задание [здесь](docs/task.md)

Реализованное консольное приложение - прототип системы для обработки соревнований по биатлону с возможностью:

- Генерации отчета гонки
- Логирования всех действий спортсменов на основе внешнего файла с событиями

## 📌 Особенности

- **Режим реального времени** - события обрабатываются с соблюдением временных меток
- **Гибкая конфигурация** - параметры гонки задаются через JSON-файл
- **Подробный отчет** - включает время кругов, штрафы, статистику стрельбы
- **Двойное логирование** - вывод логов одновременно в консоль и файл

## 🚀 Запуск

Перед запуском убедитесь, что у вас установлен golang 1.24

1. Склонируйте проект к себе на компьютер

```bash
git clone https://github.com/KTemka1234/Yadro-Impulse-2025-2.git
cd ./Yadro-Impulse-2025-2
```

2. Запустите проект

```bash
go run .

# Или используйте air live-reload
air
```

## 📂 Структура файлов

```bash
biathlon-system/
├── configs/               # Примеры конфигураций + парсинг
│   ├── config_xmpl.json
│   ├── config.go
│   └── config.json
├── docs/                  # Исходное задание
│   └── task.md
├── logger/                # Настройки логов событий
│   └── task.md
├── logs/                  # Полученные логи
│   └── task.md
├
├── main.go                # Точка входа
├── events                 # Файл событий
├── eventHandlers.go       # Логика обработки событий
├── types.go               # Используемые типы
├── go.mod
├── go.sum
└── README.md
```

## 📜 Лицензия

MIT License. Подробнее в файле [LICENSE](LICENSE).
