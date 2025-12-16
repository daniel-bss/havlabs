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

run:
	go run ./cmd

cpenv:
	cp .env .env.example

gogetproto:
	go get github.com/daniel-bss/havlabs-proto

upload_img:
	curl -X PUT \
	-H "Content-Type: image/png" \
	--upload-file ./temp.png \
	"http://localhost:9000/havlabs/images/20bf06ee-29e9-40f3-a5e3-f120176f93e8.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=media%2F20251216%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251216T071247Z&X-Amz-Expires=120&X-Amz-SignedHeaders=host&X-Amz-Signature=8c8770cfd12e8e43dc598e8c00f39759694eadd47b1cb9c12a9962b316c240b1"


.PHONY: evans sqlc new_migration migrateup migrateup1 migratedown migratedown1 run cpenv gogetproto upload_img
