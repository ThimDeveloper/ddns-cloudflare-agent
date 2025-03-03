package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"

	"github.com/thoas/go-funk"

	"gopkg.in/yaml.v3"
)

var IPAPI_URL string = "http://ip-api.com/json"
var IPIFY_URL string = "https://api.ipify.org?format=json"
var CLOUDFLARE_HOST string = "https://api.cloudflare.com/client/v4"

func check(e error) {
	if e != nil {
		slog.Error(e.Error())
		panic(e)
	}

}

func upsertDNSRecord(providerConfig *ProviderConfig, latestRouterIp string) {

	client := cloudflare.NewClient(
		option.WithAPIToken(providerConfig.DNSProvider.Cloudflare.APIToken),
		option.WithMaxRetries(3),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	records, err := client.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(providerConfig.DNSProvider.Cloudflare.ZoneID),
	})
	if err != nil {
		panic(err.Error())
	}
	relevantRecords := funk.Map(providerConfig.DNSProvider.Cloudflare.Records, func(x ConfigRecord) string {
		return x.ID
	}).([]string)

	matchingRecords := funk.Filter(records.Result, func(x dns.RecordResponse) bool {
		return funk.ContainsString(relevantRecords, x.ID)
	}).([]dns.RecordResponse)

	for _, record := range matchingRecords {
		if record.Content != latestRouterIp {
			slog.Info("IP change detected!\nOld IP: %s\nNew IP: %s\n", record.Content, latestRouterIp)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			_, err := client.DNS.Records.Edit(ctx, record.ID, dns.RecordEditParams{
				ZoneID: cloudflare.F(providerConfig.DNSProvider.Cloudflare.ZoneID),
				Record: dns.RecordParam{Content: cloudflare.F(latestRouterIp)},
			})
			if err != nil {
				slog.Error(err.Error())
			}
			slog.Info("Updated %s to IP %s\n", record.Name, latestRouterIp)
		} else {
			slog.Info("No change in IP. Skipping...")
		}
	}

}

func readConfiguration() (ProviderConfig, error) {
	configFilePath := "/etc/ddns-cloudflare-agent/config.yml"
	configFilePathFromEnv, isFilePathPresent := os.LookupEnv("OVERRIDE_CONFIGURATION_FILE_PATH")
	if isFilePathPresent {
		configFilePath = configFilePathFromEnv
	}
	var providerConfig ProviderConfig
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		slog.Error(err.Error())
		return ProviderConfig{}, err
	}
	err = yaml.Unmarshal(data, &providerConfig)
	if err != nil {
		slog.Error(err.Error())
		return ProviderConfig{}, err
	}
	return providerConfig, nil

}

func main() {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "INFO":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "WARN":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "ERROR":
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
	currentExternalIP, err := getExternalIp()
	check(err)
	configuration, err := readConfiguration()
	check(err)
	upsertDNSRecord(&configuration, currentExternalIP.IP)
}
