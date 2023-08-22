
# Backends

.PHONY: build-backends
build-backends:
	GOOS=linux go build -o build/talktome-startartconversation cmd/startartconversation/main.go
	GOOS=linux go build -o build/talktome-continueconversation cmd/continueconversation/main.go
	GOOS=linux go build -o build/talktome-getconversation cmd/getconversation/main.go
	GOOS=linux go build -o build/talktome-listconversations cmd/listconversations/main.go

.PHONY: zip-backends
zip-backends:
	zip -j build/zip/talktome-startartconversation.zip build/talktome-startartconversation
	zip -j build/zip/talktome-continueconversation.zip build/talktome-continueconversation
	zip -j build/zip/talktome-getconversation.zip build/talktome-getconversation
	zip -j build/zip/talktome-listconversations.zip build/talktome-listconversations

.PHONY: deploy-backends
deploy-backends:
	aws s3 sync build/zip/ s3://talktome-backends/

.PHONY: clean
clean:
	rm -rf build/
	mkdir build/
	mkdir build/zip/

.PHONY: release-backends
release-backends: clean build-backends zip-backends deploy-backends

.PHONY: start-mock-server
start-mock-server:
	go run cmd/mockserver/main.go

# App

.PHONY: start-app
start-app:
	cd app && npm run start

.PHONY: build-app
build-app:
	cd app && npm run build

.PHONY: deploy-app
deploy-app:
	cd app && aws s3 sync build/ s3://talktome-app/

.PHONY: release-app
release-app: build-app deploy-app
