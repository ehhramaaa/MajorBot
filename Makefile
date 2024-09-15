
build:
	docker build -t majorbot .

up:
	docker-compose up -d

down:
	docker-compose down

delete:
	docker rmi majorbot --force

.PHONY: build up down delete