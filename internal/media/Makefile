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
	"http://localhost:9000/staging/news/b4500553-391f-4283-86f6-21e23a6a852a?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=media%2F20251220%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251220T192609Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host&X-Amz-Signature=03222e78853ce11255ddad49572c65ff104298342b60756e19a53d6589a6f5a9"

backup:
	docker container run --rm --name dummybak --mount "type=bind,source=./temp/backup,destination=/backup" --mount "type=volume,source=media_media-volume,destination=/data" ubuntu:latest tar cvf /backup/backup.tar.gz /data

restore:
	docker container run --rm --name dummyres --mount "type=bind,source=./temp/backup,destination=/backup" --mount "type=volume,source=media_media-backup,destination=/data" ubuntu:latest bash -c "cd /data && tar xvf /backup/backup.tar.gz --strip 1"

.PHONY: evans sqlc new_migration migrateup migrateup1 migratedown migratedown1 run cpenv gogetproto upload_img
