package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func route53UpsertDNSRecord(providerConfig *ProviderConfig, latestRouterIp string) {
	// Load the SDK's configuration from environment and shared config, and
	// create the client with this.
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}

	client := route53.NewFromConfig(cfg)

	changes := []types.Change{}
	for _, record := range providerConfig.DNSProvider.AWS.Records {
		change := types.Change{
			Action: types.ChangeActionUpsert,
			ResourceRecordSet: &types.ResourceRecordSet{
				Name: &record.Name,
				Type: types.RRTypeA,
				ResourceRecords: []types.ResourceRecord{
					{
						Value: &latestRouterIp,
					},
				},
			},
		}
		changes = append(changes, change)
	}

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: changes,
			Comment: &providerConfig.DNSProvider.AWS.Comment,
		},
		HostedZoneId: &providerConfig.DNSProvider.AWS.HostedZoneID,
	}

	updateRes, updateErr := client.ChangeResourceRecordSets(context.Background(), params)
	if updateErr != nil {
		panic(updateErr)
	}
	slog.Info("Route53: DNS records updated", "id", updateRes.ChangeInfo.Id)

	getChangeParams := &route53.GetChangeInput{Id: updateRes.ChangeInfo.Id}

	changeRes, changeErr := client.GetChange(context.Background(), getChangeParams)
	if changeErr != nil {
		panic(changeErr)
	}
	slog.Info("Route53: DNS records status", "status", changeRes.ChangeInfo.Status)

}
