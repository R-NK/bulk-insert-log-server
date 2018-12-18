run: main.go
	go build main.go
	./main

docker/exec:
	docker-compose exec mysql bash

docker/up:
	docker-compose up

docker/up-d:
	docker-compose up -d

hey:
	hey -c 200  http://localhost:8080/hakaru\?name\=hey\?value\=100