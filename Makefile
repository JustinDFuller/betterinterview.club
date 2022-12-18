SHELL := /bin/bash

export HOST=https://localhost:8443/

server-watch: ## Run the Go server and restart on changes.
	@reflex -s -d "fancy" -r '\.(go|html)' -R 'organizations.json' -- sh -c 'clear && $(MAKE) server';

server:  ## Run the Go server.
	@go run -race ./cmd/localhost;

certificates: ## Generate SSL certificates for the HTTPS server.
	@openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr;
	@openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt;

deploy: ## Build and deploy for app engine
	@cd ./cmd/appengine; gcloud app deploy;

dispatch: ## Deploy routing rules for app engine
	@cd ./cmd/appengine; gcloud app deploy dispatch.yaml;

format:
	@npx prettier -w ./**/*.html ./**/*.css;
