
.PHONY: build
build:
	GOOS=linux go build -o build/talktome cmd/talktome/main.go

.PHONY: start-mock-server
start-mock-server:
	go run cmd/talktomeartcreate/mockserver/main.go
