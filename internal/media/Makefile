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
	curl -X PUT \
	--upload-file ./temp/a.png \
	-H "Content-Type: image/png" \ 
	"http://localhost:9000/staging/news/717fd2e1-b4f8-4ad4-80fd-3656064c69e4?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=media%2F20251218%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251218T113410Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host&X-Amz-Signature=de33b15f73a60d1e61230a07a5ec78f49277ce61f1bf763f13cd3746464930b5"

.PHONY: evans sqlc new_migration migrateup migrateup1 migratedown migratedown1 run cpenv gogetproto upload_img
