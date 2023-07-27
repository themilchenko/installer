package main

import (
	"fmt"
	"log"
	"os"

	"installer/config"
	"installer/downloader"
)

func main() {
	var pathToCfg string
	if len(os.Args) == 1 {
		pathToCfg = defaultCfgName
	} else if len(os.Args) == 2 {
		pathToCfg = os.Args[1]
	} else {
		log.Fatal("this program works with less than 2 arguments")
	}

	// Reading config
	cfg := new(config.InstallerConfig)

	err := cfg.Open(pathToCfg)
	if err != nil {
		log.Fatal("cannot open configuration file")
	}

	// Download files form server
	downloader.DownloadAllFiles(cfg.Files, cfg.ServerAddr+":"+cfg.ServerPort)

	// Start installation
	err = cfg.Installation()
	if err != nil {
		log.Fatalf("Fatal error: %s\n", err.Error())
	}

	fmt.Println("Installer has finished work successfully")
}
