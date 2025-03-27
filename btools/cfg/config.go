package cfg

import (
	"github.com/seemyown/backend-toolkit/btools/logging"
	"github.com/spf13/viper"
)

var log = logging.New(logging.Config{
	FileName: "config",
	Name:     "config",
})

func NewConfig[T struct{}](
	fileName, ext, path string,
) *T {
	var config T

	viper.SetConfigName(fileName)
	viper.SetConfigType(ext)
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Error(err, "Error reading config file")
		return nil
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Error(err, "Error unmarshalling config file")
		return nil
	}
	log.Debug("Config loaded successfully")
	return &config
}
