include .env

DB_URL=$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

evans:
	evans -r repl -p 50050 --verbose

sqlc:
	sqlc generate

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

run:
	go run ./cmd

renew_env_example:
	cp .env .env.example

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--include_imports \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--descriptor_set_out=pb/all.pb \
	proto/*.proto

gogetproto:
	go get github.com/daniel-bss/havlabs-proto@8a4bc40

.PHONY: proto run sqlc new_migration migrateup migrateup1 migratedown migratedown1 evans
