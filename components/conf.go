package components

import (
	"bytes"
	"fmt"
	"prometheus-transporter/model"

	"github.com/spf13/viper"
	"github.com/toolkits/pkg/file"
)

var (
	config   *model.Config
	Endpoint string
	Cwd      string
)

func GetConfig() *model.Config {
	return config
}

func ParseConfig(conf string) error {
	bs, err := file.ReadBytes(conf)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	viper.SetConfigType("toml")
	err = viper.ReadConfig(bytes.NewBuffer(bs))
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	setDefault()

	err = viper.Unmarshal(&config)
	if err != nil {
		return fmt.Errorf("Unmarshal %v", err)
	}

	return nil
}

func setDefault() {
	viper.SetDefault("httpPort", 8089)
	viper.SetDefault("transfers", []string{
		"127.0.0.1:5821",
	})

	viper.SetDefault("logger", map[string]interface{}{
		"level":  "INFO",
		"format": "json",
	})
}
