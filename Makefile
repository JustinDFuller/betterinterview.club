server-watch: ## Run the Go server and restart on changes.
	@reflex -s -r '\.go' -- sh -c 'clear && $(MAKE) server';

server:  ## Run the Go server.
	@go run ./cmd/server/main.go;

certificates: ## Generate SSL certificates for the HTTPS server.
	@openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr;
	@openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt;
