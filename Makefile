SHELL := /bin/bash

server-watch: ## Run the Go server and restart on changes.
	@reflex -s -r '\.go' -- sh -c 'clear && $(MAKE) server';

server:  ## Run the Go server.
	@go run ./cmd/https/main.go;

certificates: ## Generate SSL certificates for the HTTPS server.
	@openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr;
	@openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt;

deploy: ## Build and deploy for app engine
	@cd ./cmd/http; gcloud app deploy;

dispatch: ## Deploy routing rules for app engine
	@cd ./cmd/http; gcloud app deploy dispatch.yaml;

email: ## Send a test email
	@go run ./cmd/mail;
