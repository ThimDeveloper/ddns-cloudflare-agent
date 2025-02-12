package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"

	"github.com/thoas/go-funk"

	"gopkg.in/yaml.v3"
)

type ConfigRecord struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
	Type string `yaml:"type"`
}
type ProviderConfig struct {
	DNSProvider struct {
		Cloudflare struct {
			APIToken string         `yaml:"api_token"`
			ZoneID   string         `yaml:"zone_id"`
			Records  []ConfigRecord `yaml:"records"`
		} `yaml:"cloudflare"`
	} `yaml:"dns_provider"`
}

type ExternalIPResponse struct {
	IP string `json:"ip"`
}

type IPIFYResponse struct {
	IP string `json:"ip"`
}

type IPAPIResponse struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

var IPAPI_URL string = "http://ip-api.com/json"
var IPIFY_URL string = "https://api.ipify.org?format=json"
var CLOUDFLARE_HOST string = "https://api.cloudflare.com/client/v4"

func check(e error) {
	if e != nil {
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

	fmt.Printf("Fetching information for %s\n", providerConfig.DNSProvider.Cloudflare.Records)
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
			fmt.Printf("IP change detected!\nOld IP: %s\nNew IP: %s\n", record.Content, latestRouterIp)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			response, err := client.DNS.Records.Edit(ctx, record.ID, dns.RecordEditParams{
				ZoneID: cloudflare.F(providerConfig.DNSProvider.Cloudflare.ZoneID),
				Record: dns.RecordParam{Content: cloudflare.F(latestRouterIp)},
			})
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Updated %s to IP %s at %s\n", record.Name, record.Content, response.ModifiedOn.String())
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
		fmt.Println(err)
		return ProviderConfig{}, err
	}
	err = yaml.Unmarshal(data, &providerConfig)
	if err != nil {
		fmt.Println(err)
		return ProviderConfig{}, err
	}
	return providerConfig, nil

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	currentExternalIP, err := getExternalIp()
	check(err)
	configuration, err := readConfiguration()
	check(err)
	upsertDNSRecord(&configuration, currentExternalIP.IP)

}
