
.PHONY: build
build:
	GOOS=linux go build -o build/talktome cmd/talktome/main.go
	GOOS=linux go build -o build/callback cmd/resemblecallback/main.go
