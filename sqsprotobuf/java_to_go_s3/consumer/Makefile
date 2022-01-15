.PHONY: upgrade-libraries
upgrade-libraries:
	go get -u ./...
	go fmt ./...
	go mod tidy
	go mod verify

.PHONY: gen-proto
gen-proto:
	docker run --rm -v $(shell pwd):/work/output -v $(shell pwd)/../../protobuf/addressbook.proto:/work/addressbook.proto  goprotoc --proto_path=/work --go_out=/work/output addressbook.proto
