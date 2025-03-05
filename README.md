# ddns-cloudflare-agent
A simple ddns-agent for cloudflare that automatically updates DNS records to point to your home lab router IP or other public IP that might change such as AWS EC2 public IP


## Usage

### Configure the provider
Create a file named `config.yml` in a directory named `/etc/ddns-cloudflare-agent`

```yaml
dns_provider:
  cloudflare:
    api_token: "your-api-token"
    zone_id: "your-zone-id"
    records:
      - name: "domain.xyz"
        id: a123
        type: a
      - name: "*.domain.xyz"
        id: b123
        type: a
```

### Using docker-compose

#### Mount the configuration

When running the agent as a docker container is makes more sense to let the application in the container re-schedule itself on a specified interval than to reschedule the binary using cron. This way the container can be stopped and started without losing the schedule.
The also allows for easier log checks of the container.

```yaml
services:
    ddns-cloudflare-agent:
        image: ddns-cloudflare-agent:latest # or pin specific version
        environment:
            RUN_IN_DOCKER: "true"
            DOCKER_SCHEDULE_INTERVAL: 300 # override ddns update interval e.g. every 5 minutes in seconds (default 600)
        volumes:
            - /etc/ddns-cloudflare-agent:/etc/ddns-cloudflare-agent # <-- mount point in container
```

### Using the binary

### Build the binary

```bash
# Cross compile for linux, windows and mac
./build.sh
```

#### Schedule the binary using system scheduler

```bash
# Example using crontab (crontab -e)
0 0-23 * * * /path/to/ddns-cloudflare-agent
```