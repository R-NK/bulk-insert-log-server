run: main.go
	go build main.go
	./main

docker/exec:
	docker-compose exec mysql bash