.PHONY: build
build:
	cd app && env GOOS=linux GOARCH=arm64 go build -o ../bin/ddns-cloudflare-agent_arm64
	cd app && go build -o ../bin/ddns-cloudflare-agent

.PHONY: run
run:
	@cd app && go run .


IMAGE_VERSION=1.0.0

.PHONY: docker-publish
docker-publish:
	docker buildx build --platform linux/arm64,linux/amd64 -t thimlohse/ddns-cloudflare-agent:$(IMAGE_VERSION) -t thimlohse/ddns-cloudflare-agent:latest  . --push
