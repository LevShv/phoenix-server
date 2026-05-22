# Переменные
GOOS = linux
GOARCH = amd64
CGO_ENABLED = 1
BINARY = server

# Настройки сервера
REMOTE_USER = root
REMOTE_HOST = llvvv.ru
REMOTE_PATH = /home/app/

# Фантомные цели
.PHONY: all build send clean

# Цель по умолчанию
all: build

# Сборка бинарного файла
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -o $(BINARY)

# Отправка на сервер (зависит от build - сначала соберёт, потом отправит)
send: build
	scp $(BINARY) $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_PATH)

# Очистка
clean:
	rm -f $(BINARY)

# Доп. полезные цели
.PHONY: rebuild
rebuild: clean build

.PHONY: deploy
deploy: clean build send