
.PHONY: build
build:
	go build -o build/talktome cmd/talktome/main.go

.PHONY: run
run:
	go run cmd/talktome/main.go
