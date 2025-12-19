include .env

DB_URL=$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

evans:
	evans -r repl -p $(GRPC_SERVER_ADDRESS) --verbose

sqlc:
	sqlc generate

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migrateforce:
	migrate -path db/migration -database "$(DB_URL)" force 1

run:
	go run ./cmd

cpenv:
	cp .env .env.example

gogetproto:
	go get github.com/daniel-bss/havlabs-proto

upload_img:
	curl -X PUT  \
	--upload-file ./temp/a.png \
	"http://localhost:9000/staging/news/fede7aa0-6a12-4f3e-9783-f7be86638b85?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=media%2F20251219%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251219T092047Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host&X-Amz-Signature=e27f38b30bc51484c5132f93a213b15499b08e203ae6618f3e5a8be73fdaeb8b"

.PHONY: evans sqlc new_migration migrateup migrateup1 migratedown migratedown1 run cpenv gogetproto upload_img
