.PHONY: build
build:
	@cd app && go build -o ../bin/ddns-cloudflare-agent

.PHONY: run
run:
	@cd app && go run .