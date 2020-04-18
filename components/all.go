package components

import (
	"flag"
	"fmt"
)

func InitAll() error {
	// parse
	conf := flag.String("c", "./dev.conf.toml", "specify the configuration file")
	flag.Parse()

	err := ParseConfig(*conf)
	if err != nil {
		// Please noted: Error Wrap requires go 1.13
		return fmt.Errorf("init transporter error %w", err)
	}

	InitLogger(config.Logger)

	return nil
}
