package utils

import "github.com/spf13/viper"

type AppConfig struct {
	AppPort               string `mapstructure:"APP_PORT"`
	AppEnvironment        string `mapstructure:"APP_ENV"`
	AppDomain             string `mapstructure:"APP_DOMAIN"`
	AccessTokenKey        string `mapstructure:"ACCESS_TOKEN_KEY"`
	RefreshTokenKey       string `mapstructure:"REFRESH_TOKEN_KEY"`
	DatabaseName          string `mapstructure:"DATABASE_NAME"`
	BookingTopic          string `mapstructure:"BOOKING_TOPIC"`
	GoogleCredentialsPath string `mapstructure:"GOOGLE_CREDENTIALS_PATH"`
	MidtransServerKey     string `mapstructure:"MIDTRANS_SERVER_KEY"`
}

func GetConfig() (config *AppConfig, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return config, nil
}
