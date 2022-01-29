.PHONY: build
build:
	make -C ../lambda build

.PHONY: up
up: build
	aws-vault exec home -- pulumi up

.PHONY: down
down:
	aws-vault exec home -- pulumi destroy

.PHONY: upgrade-libraries
upgrade-libraries:
	@go get -u ./...
	@go fmt ./...
	@go mod tidy
	@go mod verify
