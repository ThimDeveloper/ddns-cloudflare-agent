package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/go-co-op/gocron/v2"

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
			slog.Info("IP change detected!", "old", record.Content, "new", latestRouterIp)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			_, err := client.DNS.Records.Edit(ctx, record.ID, dns.RecordEditParams{
				ZoneID: cloudflare.F(providerConfig.DNSProvider.Cloudflare.ZoneID),
				Record: dns.RecordParam{Content: cloudflare.F(latestRouterIp)},
			})
			if err != nil {
				slog.Error(err.Error())
			}
			slog.Info("Updated IP", "old", record.Name, "new", latestRouterIp)
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

func runDdns() {
	currentExternalIP, err := getExternalIp()
	check(err)
	configuration, err := readConfiguration()
	check(err)
	upsertDNSRecord(&configuration, currentExternalIP.IP)
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
	runInDocker := os.Getenv("RUN_IN_DOCKER") == "true"
	if runInDocker {
		slog.Debug("Running in Docker")
		dockerScheduleInterval, isDockerScheduleIntervalPresent := os.LookupEnv("DOCKER_SCHEDULE_INTERVAL")

		var scheduleInterval time.Duration = time.Duration(600) * time.Second
		if isDockerScheduleIntervalPresent {
			overrideInterval, err := strconv.Atoi(dockerScheduleInterval)
			check(err)
			slog.Debug("Setting schedule interval", "scheduleIntervalSeconds", time.Duration(overrideInterval)*time.Second)
			scheduleInterval = time.Duration(overrideInterval) * time.Second
		} else {
			slog.Debug("Using default schedule interval", "scheduleIntervalSeconds", "600s")
		}

		s, err := gocron.NewScheduler()
		if err != nil {
			check(err)
		}

		j, err := s.NewJob(
			gocron.DurationJob(
				scheduleInterval,
			),
			gocron.NewTask(runDdns),
		)
		if err != nil {
			check(err)
		}
		slog.Debug(j.ID().String())
		s.Start()

		select {} // block forever until the program is terminated
	} else {
		slog.Debug("Running binary outside of Docker")
		runDdns()
	}
}
