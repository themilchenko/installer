package main

import (
	"fmt"
	"log"

	"installer/config"
)

func main() {
	cfg := new(config.InstallerConfig)
	err := cfg.GenerateConfig()
	if err != nil {
		log.Fatalf("Cannot create config: %s", err.Error())
	}
	fmt.Println("\nConfig was successfully created with your variables")
}
