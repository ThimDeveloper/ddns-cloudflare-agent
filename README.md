# ddns-cloudflare-agent
A simple ddns-agent for cloudflare that automatically updates DNS records to point to your home lab router IP or other public IP that might change such as AWS EC2 public IP


## Usage

### Configure the provider
Create a file named `config.yml` in `/etc/ddns-cloudflare-agent`

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
        image: ddns-cloudflare-agent:latest # or pin specific version
        volumes:
            - /etc/ddns-cloudflare-agent:/etc/ddns-cloudflare-agent # <-- mount point in container
```

## Schedule the container to run on fixed schedule

Setup small shell script to invoke docker for easier reference in crontab:
```bash
# ddns-cloudflare-agent in /usr/local/bin
echo "Invoking ddns-cloudflare-agent"
docker compose -f <path-to-docker-compose> down
docker compose -f <path-to-docker-compose> up -d --remove-orphans
```

Make script executable:
```bash
chmod +x /usr/local/bin/ddns-cloudflare-agent
```

```bash
# Example using crontab (crontab -e)
0 0-23 * * * ddns-cloudflare-agent
```