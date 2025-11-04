package configs

type Provider struct {
	ProviderType string `yaml:"provider_type"`
	NetworkName  string `yaml:"network_name"`
	BaseURL      string `yaml:"base_url"`
	ApiKey       string `yaml:"api_key"`
}
