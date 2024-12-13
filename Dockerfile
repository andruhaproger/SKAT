# Используем официальный образ Go
FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Компилируем приложение
RUN go build -o main .

# Указываем команду для запуска приложения
CMD ["./main"]
