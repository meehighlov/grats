# grats

## Запуск

1. создаем .env файл и добавляем его в /cmd/grats (пример .env.example)
2. запускаем миграции
   - устанавливаем goose, например: brew install goose
   - из каталога /cmd/grats запускаем команду
   ```shell
   goose -dir=../../migrations sqlite3 grats.db up
   ```
3. запускаем бота из каталога /cmd/grats:
   ```shell
   go run main.go
   ```
