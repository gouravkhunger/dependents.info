API_DIR = api
ACTION_DIR = action

clean:
	@cd $(API_DIR) && go clean
	@cd $(ACTION_DIR) && rm -rf node_modules dist

install:
	@cd $(API_DIR) && go mod tidy
	@cd $(ACTION_DIR) && npm install -s

api-dev:
	@cd $(API_DIR) && go run main.go

api-build: install
	@cd $(API_DIR) && go build -o bin/api main.go -ldflags="-s -w"

api-test:
	@cd $(API_DIR) && go test ./...

action-build: install
	@cd $(ACTION_DIR) && npm run build

action-local:
	@cd $(ACTION_DIR) && npm run local-action

action-test:
	@cd $(ACTION_DIR) && npm run test

.PHONY: clean install api-dev api-build api-test action-build action-local action-test
