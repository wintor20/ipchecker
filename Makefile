.PHONY: up down build mod

build:
	docker-compose build

mod:
	rm go.mod
	go mod init ipchecker

up:
	docker-compose -p ipchecker --log-level ERROR up -d

down:
	docker-compose -p ipchecker down -v

logs:
	docker-compose -p ipchecker logs -f checker


