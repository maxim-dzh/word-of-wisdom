test:
	go test -v ./...

start-server:
	go run cmd/server/main.go

start-client:
	go run cmd/client/main.go

build-and-start-docker:
	docker compose up -d --force-recreate --build server --build client
	docker logs word-of-wisdom-server
	docker logs word-of-wisdom-client

start-docker:
	docker compose up -d --force-recreate
	docker logs word-of-wisdom-server
	docker logs word-of-wisdom-client

lint:
	golangci-lint -v run ./...
