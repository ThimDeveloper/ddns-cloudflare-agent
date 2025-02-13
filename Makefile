.PHONY: build
build:
	@cd app && env GOOS=linux GOARCH=arm64 go build -o ../bin/ddns-cloudflare-agent_arm64
	@cd app && go build -o ../bin/ddns-cloudflare-agent

.PHONY: run
run:
	@cd app && go run .
