# ddns-cloudflare-agent
A simple ddns-agent for cloudflare that automatically updates DNS records to point to your home lab router IP or other public IP that might change such as AWS EC2 public IP


## Usage

### Configure the provider
```yaml
dns_provider:
  cloudflare:
    api_token: "your-api-token"
    zone_id: "your-zone-id"
    records:
      - name: domain.xyz
        id: a123
        type: a
      - name: "example.codedungeon.xyz"
        id: b123
        type: a
```

### Mount the configuration
```yaml
services:
    ddns-cloudflare-agent:
        container_name: ddns-cloudflare-agent
        image: ddns-cloudflare-agent:latest # or pin specific version
        volumes:
            - ./config:/etc/ddns-cloudflare-agent # <-- mount point in container
```
