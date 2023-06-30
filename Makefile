
.PHONY: build
build:
	go build -o build/talktome cmd/talktome/main.go
	go build -o build/callback cmd/resemblecallback/main.go
