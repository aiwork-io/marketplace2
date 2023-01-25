package internal

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

const CTX_CONFIGS string = "internal.configs"

type Configs struct {
	Version                        string `mapstructure:"APP_VERSION"`
	Secret                         string `mapstructure:"APP_SECRET"`
	DbUri                          string `mapstructure:"APP_DB_URI"`
	Region                         string `mapstructure:"APP_REGION"`
	StorageBucket                  string `mapstructure:"APP_STORAGE_BUCKET"`
	TxnRate                        int    `mapstructure:"APP_TXN_RATE"`
	TxnValidatorEndpoint           string `mapstructure:"APP_TXN_VALIDTOR_ENDPOPINT"`
	TxnValidatorAuthtoken          string `mapstructure:"APP_TXN_VALIDTOR_AUTH_TOKEN"`
	TxnValidatorPaymentToken       string `mapstructure:"APP_TXN_VALIDTOR_PAYMENT_TOKEN"`
	TxnValidatorDestAddress        string `mapstructure:"APP_TXN_VALIDTOR_DEST_ADDRESS"`
	AuthVerifictionResetEndpoint   string `mapstructure:"APP_AUTH_VERIFICATION_RESET_ENDPOINT"`
	AuthVerifictionAccountEnable   int    `mapstructure:"APP_AUTH_VERIFICATION_ACCOUNT_ENABLE"`
	AuthVerifictionAccountEndpoint string `mapstructure:"APP_AUTH_VERIFICATION_ACCOUNT_ENDPOINT"`

	TestEmail string `mapstructure:"TEST_EMAIL"`
}

func NewConfigs(provider *viper.Viper) *Configs {
	configs := Configs{}

	// common configs
	if err := provider.Unmarshal(&configs); err != nil {
		panic(err)
	}

	if configs.Secret == "" {
		panic(errors.New("APP_SECRET was not set"))
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

	provider.SetDefault("APP_VERSION", getVersion())
	provider.SetDefault("APP_SECRET", "changemenow")
	provider.SetDefault("APP_DB_URI", "postgresql://aiworkmarketplace:changemenow@127.0.0.1:5432/aiworkmarketplace?sslmode=disable")
	provider.SetDefault("APP_REGION", "ap-southeast-1")
	provider.SetDefault("APP_STORAGE_BUCKET", "aiworkmarketplace-private")
	provider.SetDefault("APP_TXN_RATE", 3)
	provider.SetDefault("APP_AUTH_VERIFICATION_ACCOUNT_ENABLE", 1)

	provider.AutomaticEnv()

	return provider
}

func getVersion() string {
	if body, err := os.ReadFile(".version"); err == nil {
		return string(body)
	}

	return "22.2.2"
}
