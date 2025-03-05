.PHONY: run
run:
	@cd app && go run .


IMAGE_VERSION=2.0.0

.PHONY: docker-publish
docker-publish:
	docker buildx build --platform linux/arm64,linux/amd64 -t thimlohse/ddns-cloudflare-agent:$(IMAGE_VERSION) -t thimlohse/ddns-cloudflare-agent:latest  . --push
