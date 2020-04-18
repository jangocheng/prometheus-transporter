package main

import (
	"prometheus-transporter/components"
	"prometheus-transporter/receiver"
)

func main() {
	components.InitComponents()
	logger := components.GetLogger()
	conf := components.GetConfig()
	receiver.Start(logger, conf)
}
