package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/thoas/go-funk"
)

var CLOUDFLARE_HOST string = "https://api.cloudflare.com/client/v4"

func cloudflareUpsertDNSRecord(providerConfig *ProviderConfig, latestRouterIp string) {

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
			slog.Info("Cloudflare: IP change detected!", "old", record.Content, "new", latestRouterIp)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			_, err := client.DNS.Records.Edit(ctx, record.ID, dns.RecordEditParams{
				ZoneID: cloudflare.F(providerConfig.DNSProvider.Cloudflare.ZoneID),
				Record: dns.RecordParam{Content: cloudflare.F(latestRouterIp)},
			})
			if err != nil {
				slog.Error(err.Error())
			}
			slog.Info("Cloudflare: Updated IP", "old", record.Name, "new", latestRouterIp)
		} else {
			slog.Info("Cloudflare: No change in IP. Skipping...")
		}
	}

}
