package components

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"github.com/toolkits/pkg/file"
)

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type ConfToml struct {
	HTTPPort  int       `toml:"httpPort"`
	Transfers []string  `toml:"transfers"`
	Logger    LogConfig `toml:"logger"`
}

var (
	config   *ConfToml
	lock     = new(sync.RWMutex)
	Endpoint string
	Cwd      string
)

func GetConfig() *ConfToml {
	return config
}

func ParseConfig(conf string) error {
	bs, err := file.ReadBytes(conf)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	lock.Lock()
	defer lock.Unlock()

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
