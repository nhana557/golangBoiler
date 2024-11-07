start: 
	docker compose up -d --build

dependencies-up: 
	docker compose up --scale api=0

dependencies-down: 
	docker compose down mongo redis rabbitmq

server:
	./mainServer

build:
	go build -o mainServer ./main.go

.PHONEY: chmod +x manage.sh