API_DIR = api
WWW_DIR = www
ACTION_DIR = action

clean:
	@cd $(API_DIR) && go clean
	@cd $(WWW_DIR) && rm -rf node_modules dist
	@cd $(ACTION_DIR) && rm -rf node_modules dist

install:
	@cd $(API_DIR) && go mod tidy
	@cd $(WWW_DIR) && npm install -s
	@cd $(ACTION_DIR) && npm install -s

www-dev:
	@cd $(WWW_DIR) && npm run dev

www-build:
	@cd $(WWW_DIR) && npm run build

api-fmt:
	@cd $(API_DIR) && gofmt -w .

api-dev:
	@cd $(WWW_DIR) && npm run build
	@mkdir -p $(API_DIR)/static
	@cp -r $(WWW_DIR)/dist/* $(API_DIR)/static/
	@cd $(API_DIR) && go run main.go

api-build:
	@cd $(WWW_DIR) && npm run build
	@mkdir -p $(API_DIR)/static
	@cp -r $(WWW_DIR)/dist/* $(API_DIR)/static/
	@cd $(API_DIR) && go build -o bin/api main.go

api-test:
	@cd $(API_DIR) && go test ./...

action-build:
	@cd $(ACTION_DIR) && npm run build

action-local:
	@cd $(ACTION_DIR) && npm run local-action

action-test:
	@cd $(ACTION_DIR) && npm run test
