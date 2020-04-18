package main

import (
	"transporter/components"
)

func main() {
	components.InitAll()
	logger := components.GetLogger()
}
