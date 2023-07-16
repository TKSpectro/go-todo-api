.PHONY: start
start:
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down

.PHONY: build
build:
	go build -o ./tmp/make_build ./main.go

.PHONY: run
run:
	go run ./main.go

.PHONY: build-run
build-run: build run

.PHONY: migrate-gen
migrate-gen:
	atlas migrate diff --env gorm $(name)

.PHONY: migrate-new
migrate-new:
	atlas migrate new --env gorm $(name)

.PHONY: migrate-up
migrate-up:
	atlas migrate apply --env local

# See: https://atlasgo.io/versioned/apply#down-migrations
.PHONY: migrate-down
migrate-down:
	atlas schema apply --env local --to "file://migrations?version=$(version)&format=golang-migrate" --exclude "atlas_schema_revisions"
	atlas migrate set --env local $(version)