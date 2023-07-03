
.PHONY: build
build:
	GOOS=linux go build -o build/talktome cmd/talktome/main.go
