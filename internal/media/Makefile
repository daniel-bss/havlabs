include .env

DB_URL=$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

evans:
	evans -r repl -p 50050 --verbose

run:
	go run ./cmd

cpenv:
	cp .env .env.example

gogetproto:
	go get github.com/daniel-bss/havlabs-proto


.PHONY: cpenv run evans
