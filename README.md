## Todo лист с авторизацией на Golang
Запуск
```sh
go run ./cmd/app/main.go
```
Отчистка зависимостей
```sh
go mod tidy
```
Билд
```sh
go build -o todo.exe ./cmd/app
```
Удалить бинарный файл
```sh
del todo.exe
```
Запуск в docker
```sh
docker-compose up -d
```