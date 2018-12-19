run: main.go
	go build main.go
	./main

docker/exec:
	docker-compose exec mysql bash

docker/up:
	docker-compose up

docker/up-d:
	docker-compose up -d

docker/down:
	docker-compose down

docker/restart:
	docker-compose restart

wrk:
	wrk -t12 -c1000 -d10s  http://localhost:8081/hakaru?name=hey?value=100