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
	"http://localhost:9000/staging/news/cf3db533-8e90-46af-80f9-9588eca839f7?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=media%2F20251220%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251220T125420Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host&X-Amz-Signature=96f919203a435ef960cbed7eebf6a5d891ce5b7a29593e48c26cdb52e94d2426"

backup:
	docker container run --rm --name dummybak --mount "type=bind,source=./temp/backup,destination=/backup" --mount "type=volume,source=media_media-volume,destination=/data" ubuntu:latest tar cvf /backup/backup.tar.gz /data

restore:
	docker container run --rm --name dummyres --mount "type=bind,source=./temp/backup,destination=/backup" --mount "type=volume,source=media_media-backup,destination=/data" ubuntu:latest bash -c "cd /data && tar xvf /backup/backup.tar.gz --strip 1"

.PHONY: evans sqlc new_migration migrateup migrateup1 migratedown migratedown1 run cpenv gogetproto upload_img
