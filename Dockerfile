FROM golang:1.24.0-alpine3.21@sha256:5429efb7de864db15bd99b91b67608d52f97945837c7f6f7d1b779f9bfe46281 AS base


FROM base AS builder
# Bulilder
WORKDIR /app
COPY app .
RUN mkdir -p /app/bin && go build -o /app/bin/ddns-cloudflare-agent


# Runner
FROM base AS runner
# Create configuration directory
RUN mkdir -p /etc/ddns-cloudflare-agent
COPY --from=builder /app/bin/ddns-cloudflare-agent /usr/local/bin/ddns-cloudflare-agent
RUN chmod +x /usr/local/bin/ddns-cloudflare-agent

CMD ["ddns-cloudflare-agent"]