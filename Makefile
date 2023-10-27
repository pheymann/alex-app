# All
.PHONY: run-contract-tests
run-contract-tests: run-backend-contract-tests run-app-contract-tests


# Backends

.PHONY: build-backends
build-backends:
	GOOS=linux go build -o build/talktome-startartconversation cmd/startartconversation/main.go
	GOOS=linux go build -o build/talktome-continueconversation cmd/continueconversation/main.go
	GOOS=linux go build -o build/talktome-getconversation cmd/getconversation/main.go
	GOOS=linux go build -o build/talktome-listconversations cmd/listconversations/main.go
	GOOS=linux go build -o build/talktome-applogs cmd/applogs/main.go
	GOOS=linux go build -o build/talktome-pollassistantresponse cmd/pollassistantresponse/main.go
	GOOS=linux go build -o build/talktome-assistant cmd/assistant/main.go

.PHONY: zip-backends
zip-backends:
	zip -j build/zip/talktome-startartconversation.zip build/talktome-startartconversation
	zip -j build/zip/talktome-continueconversation.zip build/talktome-continueconversation
	zip -j build/zip/talktome-getconversation.zip build/talktome-getconversation
	zip -j build/zip/talktome-listconversations.zip build/talktome-listconversations
	zip -j build/zip/talktome-applogs.zip build/talktome-applogs
	zip -j build/zip/talktome-pollassistantresponse.zip build/talktome-pollassistantresponse
	zip -j build/zip/talktome-assistant.zip build/talktome-assistant

.PHONY: deploy-backends
deploy-backends:
	aws s3 sync build/zip/ s3://talktome-backends/
	aws lambda update-function-code --function-name talktome-startartconversation --s3-bucket talktome-backends --s3-key talktome-startartconversation.zip
	aws lambda update-function-code --function-name talktome-continueconversation --s3-bucket talktome-backends --s3-key talktome-continueconversation.zip
	aws lambda update-function-code --function-name talktome-getconversation --s3-bucket talktome-backends --s3-key talktome-getconversation.zip
	aws lambda update-function-code --function-name talktome-listconversations --s3-bucket talktome-backends --s3-key talktome-listconversations.zip
	aws lambda update-function-code --function-name talktome-applogs --s3-bucket talktome-backends --s3-key talktome-applogs.zip
	aws lambda update-function-code --function-name talktome-pollassistantresponse --s3-bucket talktome-backends --s3-key talktome-pollassistantresponse.zip
	aws lambda update-function-code --function-name talktome-assistant --s3-bucket talktome-backends --s3-key talktome-assistant.zip

.PHONY: clean
clean:
	rm -rf build/
	mkdir build/
	mkdir build/zip/

.PHONY: release-backends
release-backends: clean build-backends zip-backends deploy-backends

.PHONY: start-mock-server
start-mock-server:
	go run cmd/testserver/main.go --mode="mock"

.PHONY: start-prod-server
start-prod-server:
	go run cmd/testserver/main.go --mode="prod"

.PHONY: run-backend-contract-tests
run-backend-contract-tests:
	go test ./internal/intergrationtest/cdc/...

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

.PHONY: run-app-contract-tests
run-app-contract-tests:
	cd app && CI=true npm run test
