
.PHONY: build
build:
	GOOS=linux go build -o build/talktome-startartconversation cmd/startartconversation/main.go
	GOOS=linux go build -o build/talktome-continueconversation cmd/continueconversation/main.go
	GOOS=linux go build -o build/talktome-getconversation cmd/getconversation/main.go
	GOOS=linux go build -o build/talktome-listconversations cmd/listconversations/main.go

.PHONY: zip
zip:
	zip -j build/talktome-startartconversation.zip build/talktome-startartconversation
	zip -j build/talktome-continueconversation.zip build/talktome-continueconversation
	zip -j build/talktome-getconversation.zip build/talktome-getconversation
	zip -j build/talktome-listconversations.zip build/talktome-listconversations

.PHONY: clean
clean:
	rm -rf build/
	mkdir build/

.PHONY: build-zip
build-zip: clean build zip

.PHONY: start-mock-server
start-mock-server:
	go run cmd/mockserver/main.go
