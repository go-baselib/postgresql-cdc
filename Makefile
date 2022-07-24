.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose kill && docker-compose rm

.PHONY: test
test:
	go test -v ./...

.PHONY: sample
sample:
	docker-compose exec -T postgres psql -U postgres < ./cmd/sample/db/init.sql
	docker-compose exec -T postgres psql -U postgres < ./cmd/sample/db/slot.sql
	cd ./cmd/sample && go run main.go
