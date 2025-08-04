# Screen Arch

Утилита для автоматической сортировки скриншотов и записей экрана на macOS.

## Описание

Сканирует рабочий стол (`~/Desktop`) и перемещает скриншоты и записи экрана в `~/Documents/Captures/` с сортировкой по типу, году и месяцу.

## Поддерживаемые файлы

- **Скриншоты**: `.png`, `.jpeg`, `.jpg` → `~/Documents/Captures/Screenshots/YYYY/MM/`
- **Записи экрана**: `.mov`, `.mp4` → `~/Documents/Captures/Screen Recordings/YYYY/MM/`

Поддерживаются имена файлов в форматах:

- `ScreenShot YYYY-MM-DD at HH.MM.SS.png`
- `Screenshot YYYY-MM-DD at HH.MM.SS.png`
- `Screen Recording YYYY-MM-DD at HH.MM.SS.mov`
- Аналоги с русскими названиями (`ScreenShot YYYY-MM-DD в HH.MM.SS.png`)

## Установка

```bash
go build -o screen-arch .
```

## Настройка фонового агента (launchd)

1. Скопируйте `.plist` в директорию LaunchAgents:

   ```bash
   cp screen-arch.plist ~/Library/LaunchAgents/
   ```

2. Загрузите агент:

   ```bash
   launchctl load ~/Library/LaunchAgents/screen-arch.plist
   ```

Агент запускается однократно при входе в систему, сортирует файлы и завершается.
