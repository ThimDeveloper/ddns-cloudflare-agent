package main

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
