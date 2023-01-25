package internal

import (
	"github.com/spf13/viper"
)

const CTX_CONFIGS string = "internal.configs"

type Configs struct {
	JobSchedule    string `mapstructure:"JOB_SCHEDULE"`
	EngineEndpoint string `mapstructure:"ENGINE_ENDPOINT"`
	ApiEndpoint    string `mapstructure:"API_ENDPOINT"`
	AuthUsername   string `mapstructure:"AUTH_USERNAME"`
	AuthPassword   string `mapstructure:"AUTH_PASSWORD"`
}

func NewConfigs(provider *viper.Viper) *Configs {
	configs := Configs{}

	// common configs
	if err := provider.Unmarshal(&configs); err != nil {
		panic(err)
	}

	return &configs
}

func NewConfigProvider(dirs ...string) *viper.Viper {
	provider := viper.New()

	provider.SetConfigName("configs")
	provider.SetConfigType("env")
	for _, dir := range dirs {
		provider.AddConfigPath(dir)
		// ignore merge error because that always be not found configs error
		_ = provider.MergeInConfig()
	}

	provider.SetDefault("API_ENDPOINT", "https://aiworkmarketplace.kubeplusplus.com")
	provider.SetDefault("JOB_SCHEDULE", "@every 1m")

	provider.AutomaticEnv()

	return provider
}
