package setting

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Database struct {
	Host              string `mapstructure:"host"`
	Port              string `mapstructure:"port"`
	User              string `mapstructure:"user"`
	Passsword         string `mapstructure:"password"`
	Name              string `mapstructure:"name"`
	MaxOpenConnection int32  `mapstructure:"max_open_connection"`
	MaxIdleConnection int32  `mapstructure:"max_idle_connection"`
}

type GoogleConnection struct {
	ApiKey string `mapstructure:"api_key"`
}

type Configuration struct {
	Database         Database         `mapstructure:"database"`
	GoogleConnection GoogleConnection `mapstructure:"google_connection"`
}

func LoadConfiguration(filePath, fileName, extension string, config *Configuration) error {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)
	viper.AddConfigPath(filePath)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.Unmarshal(config)
	return nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
