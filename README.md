<p align="center">
  <img src="build/appicon.png" alt="Knit Logo" width="128" height="128">
</p>

<h1 align="center">Knit</h1>

<p align="center">
  <strong>🎬 Элегантное десктопное приложение для поиска торрентов фильмов</strong>
</p>

<p align="center">
  <a href="https://github.com/Cheviiot/Knit/releases"><img src="https://img.shields.io/github/v/release/Cheviiot/Knit?style=flat-square&color=blue" alt="Release"></a>
  <a href="https://github.com/Cheviiot/Knit/blob/main/LICENSE"><img src="https://img.shields.io/github/license/Cheviiot/Knit?style=flat-square" alt="License"></a>
  <a href="https://github.com/Cheviiot/Knit/actions"><img src="https://img.shields.io/github/actions/workflow/status/Cheviiot/Knit/release.yml?style=flat-square" alt="Build"></a>
</p>

<p align="center">
  <a href="#-особенности">Особенности</a> •
  <a href="#-установка">Установка</a> •
  <a href="#%EF%B8%8F-настройка">Настройка</a> •
  <a href="#-сборка-из-исходников">Сборка</a> •
  <a href="#-лицензия">Лицензия</a>
</p>

---

## ✨ Особенности

- 🔍 **Поиск фильмов** — интеграция с TMDB для поиска и просмотра информации о фильмах
- 🧲 **Поиск торрентов** — автоматический поиск через публичные Jackett серверы (без настройки!)
- ⬇️ **Загрузка в один клик** — отправка magnet-ссылок в Free Download Manager
- 🎨 **Современный интерфейс** — минималистичный дизайн с тёмной темой
- 🖼️ **Красивые карточки** — постеры, рейтинги, описания фильмов
- 📱 **Адаптивный дизайн** — корректное отображение на любом размере окна
- 🚀 **Нативная производительность** — построен на Go + Wails v3

## 🛠️ Технологии

| Backend | Frontend | Интеграции |
|---------|----------|------------|
| Go 1.24 | Vue 3 | TMDB API |
| Wails v3 | Vite | Публичные Jackett серверы |
| | Tailwind CSS 4 | Free Download Manager |

## 📦 Установка

### 🐧 Linux (один клик)

```bash
curl -fsSL https://raw.githubusercontent.com/Cheviiot/Knit/main/install.sh | bash
```

или с wget:

```bash
wget -qO- https://raw.githubusercontent.com/Cheviiot/Knit/main/install.sh | bash
```

<details>
<summary>🎨 Иконка не отображается?</summary>

На некоторых дистрибутивах (ALT Linux, старые версии GNOME) иконка может не появиться автоматически. Выполните:

```bash
curl -fsSL https://raw.githubusercontent.com/Cheviiot/Knit/main/install.sh | bash -s -- --icons
```

</details>

<details>
<summary>🗑️ Удаление</summary>

```bash
curl -fsSL https://raw.githubusercontent.com/Cheviiot/Knit/main/install.sh | bash -s -- --uninstall
```

</details>

### Готовые сборки

Скачайте последнюю версию для вашей платформы в разделе [Releases](https://github.com/Cheviiot/Knit/releases):

| Платформа | Файл |
|-----------|------|
| **Linux** | `knit-x86_64.AppImage` |
| **Windows** | `knit-windows-amd64.exe` |
| **macOS Intel** | `Knit-amd64.dmg` |
| **macOS Apple Silicon** | `Knit-arm64.dmg` |

### Требования

- [Free Download Manager](https://www.freedownloadmanager.org/) — для загрузки

## ⚙️ Настройка

При первом запуске откройте настройки (⚙️) и укажите:

| Параметр | Описание | Пример |
|----------|----------|--------|
| **Jackett сервер** | Публичный сервер для поиска торрентов | `jacred.xyz`, `jac-red.ru`, `jac.red` |
| **TMDB прокси** | Прокси для получения информации о фильмах | `apn_render` |

> 💡 Приложение использует **публичные Jackett серверы** — не требуется установка и настройка собственного Jackett!

## 🔨 Сборка из исходников

### Требования

- [Go 1.24+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Task](https://taskfile.dev/) (task runner)
- Wails CLI v3

### Установка Wails CLI v3

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

### Установка Task

```bash
# macOS
brew install go-task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

# Windows (scoop)
scoop install task
```

### Разработка

```bash
wails3 dev
```

### Сборка

```bash
# Сборка для текущей ОС
wails3 build

# Или через Task
task build
```

### Кроссплатформенная сборка через Docker

```bash
# Настройка Docker образа (один раз)
task setup:docker

# Linux из любой ОС
task linux:build

# Windows из Linux/macOS
task windows:build

# macOS (только из macOS)
task darwin:build
```

### Сборка пакетов для Linux

```bash
# AppImage
task linux:appimage

# DEB/RPM пакеты
task linux:package
```

## 🎯 Поддерживаемые платформы

| Платформа | Архитектуры | Статус |
|-----------|-------------|--------|
| Linux | amd64 | ✅ Полная поддержка |
| Windows | amd64 | ✅ Полная поддержка |
| macOS | amd64, arm64 | ✅ Полная поддержка |

> 📱 **Android/iOS**: Поддержка мобильных платформ появится когда Wails v3 выйдет из alpha-версии.

## 📁 Структура проекта

```
Knit/
├── main.go                # Точка входа Wails v3
├── app.go                 # Основная логика приложения
├── torrent.go             # Работа с FDM
├── Taskfile.yml           # Система сборки
├── frontend/
│   └── src/
│       ├── App.vue        # Vue компонент интерфейса
│       ├── style.css      # Глобальные стили
│       └── assets/
│           └── Knit.png   # Иконка для UI
├── build/
│   ├── config.yml         # Конфигурация Wails v3
│   ├── appicon.png        # Иконка приложения
│   ├── darwin/            # Конфигурация macOS
│   ├── windows/           # Конфигурация Windows
│   └── linux/             # Конфигурация Linux
└── .github/
    └── workflows/
        └── release.yml    # CI/CD сборка
```

## 🤝 Вклад в проект

Pull requests приветствуются! Для крупных изменений сначала откройте [issue](https://github.com/Cheviiot/Knit/issues).

## 📄 Лицензия

[MIT](LICENSE)

---

<p align="center">
  Сделано с ❤️ by <a href="https://github.com/Cheviiot">Cheviiot</a>
</p>
