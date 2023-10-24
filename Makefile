################################################################################
########################## Docker Compose shortcuts ############################
################################################################################

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down

################################################################################
################################ Go shortcuts ##################################
################################################################################

.PHONY: build
build:
	go build -o ./tmp/make_build ./main.go

.PHONY: run
run:
	go run ./main.go

.PHONY: build-run
build-run: build run

################################################################################
################################### Ginkgo #####################################
################################################################################

.PHONY: test
test:
	GTA_ROOT_PATH=$(CURDIR) IS_TEST=true ginkgo -r $(ARGS)

.PHONY: test-v
test-v:
	GTA_ROOT_PATH=$(CURDIR) IS_TEST=true ginkgo -r -v $(ARGS)

.PHONY: test-ci
test-ci:
	GTA_ROOT_PATH=$(CURDIR) IS_TEST=true ginkgo -r --race --trace --json-report=report.json

.PHONY: test-coverage
test-coverage:
	GTA_ROOT_PATH=$(CURDIR) IS_TEST=true ginkgo -r --cover
	go tool cover -html=coverprofile.out

################################################################################
################################ Swagger Docs ##################################
################################################################################

# Will generate swagger docs in ./api
.PHONY: swag
swag:
	swag init -o ./api

################################################################################
######################## Atlas shortcuts (Migrations) ##########################
################################################################################

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