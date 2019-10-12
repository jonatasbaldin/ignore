.DEFAULT_GOAL := test

# Clean up
clean:
	@rm -fR ./cover.* ./build/ ./vendor/
.PHONY: clean

# Creates folders and download dependencies
configure:
	@mkdir -p ./build
	@GO111MODULE=on go mod download
.PHONY: configure

# Run tests and generates html coverage file
cover: test
	@go tool cover -html=./cover.out -o ./cover.html
	@rm ./cover.out
.PHONY: cover

# Run tests
test:
	@go test -v -covermode=atomic -coverprofile=cover.out ./...
.PHONY: test
