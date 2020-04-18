package components

import (
	"flag"
	"fmt"
)

func InitComponents() error {
	// parse
	conf := flag.String("c", "./dev.conf.toml", "specify the configuration file")
	flag.Parse()

	err := ParseConfig(*conf)
	if err != nil {
		// Please note: Error Wrap requires go 1.13
		return fmt.Errorf("init transporter error %w", err)
	}

	InitLogger(config.Logger)

	return nil
}
