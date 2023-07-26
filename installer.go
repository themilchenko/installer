package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"

	"agent/installer/downloader"
)

const (
	defaultCfgContent = `services:
 agent:
   name: %s
   display_name: %s
   description: %s
   host: %s
   port: %s

 observer:
   name: %s
   display_name: %s
   description: %s
   host: %s
   port: %s`

	defaultCfgName = "config.yaml"

	hostAddrFlg   = "-HostAddress="
	portFlg       = "-Port="
	srvDscFlg     = "-ServiceDescription="
	srvDspNameFlg = "-ServiceDisplayName="
	srvNameFlg    = "-ServiceName="
)

type InstallerConfig struct {
	Services struct {
		Agent struct {
			Name        string `yaml:"name"`
			DisplayName string `yaml:"display_name"`
			Description string `yaml:"description"`
			Host        string `yaml:"host"`
			Port        string `yaml:"port"`
		} `yaml:"agent"`

		Observer struct {
			Name        string `yaml:"name"`
			DisplayName string `yaml:"display_name"`
			Description string `yaml:"description"`
			Host        string `yaml:"host"`
			Port        string `yaml:"port"`
		} `yaml:"observer"`
	} `yaml:"services"`

	Installer struct {
		PathToCfg string `yaml:"path_to_config"`
		PathToBin string `yaml:"path_to_bin"`
	} `yaml:"installer"`

	Files      []string `yaml:"files"`
	ServerAddr string   `yaml:"address"`
	ServerPort string   `yaml:"port"`
}

func (c *InstallerConfig) Open(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(c); err != nil {
		return err
	}

	return nil
}

func getCorrectPath(path string) string {
	if path[len(path)-1] != '/' {
		return path + "/"
	}
	return path
}

// This function init configuration folder and create config file with default settings
func (c *InstallerConfig) installation() error {
	// Creating folder and config
	if _, err := os.Stat(c.Installer.PathToCfg + "agent.base/"); err != nil {
		if err := os.Mkdir(c.Installer.PathToCfg+"agent.base/", os.ModePerm); err != nil {
			return err
		}
	}

	c.Installer.PathToCfg = getCorrectPath(c.Installer.PathToCfg)
	c.Installer.PathToBin = getCorrectPath(c.Installer.PathToBin)

	file, err := os.Create(c.Installer.PathToCfg + "agent.base/config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	// Writing content in config
	content := fmt.Sprintf(
		defaultCfgContent,
		c.Services.Agent.Name,
		c.Services.Agent.DisplayName,
		c.Services.Agent.Description,
		c.Services.Agent.Host,
		c.Services.Agent.Port,

		c.Services.Observer.Name,
		c.Services.Observer.DisplayName,
		c.Services.Observer.Description,
		c.Services.Observer.Host,
		c.Services.Observer.Port,
	)
	_, err = fmt.Fprint(file, content)
	if err != nil {
		return err
	}

	// Building binary files for agent and observer
	// _, err = exec.Command("go", "build", "-o", c.Installer.PathToBin, "../agent/service_agent.go").
	// 	Output()
	// if err != nil {
	// 	return err
	// }
	// _, err = exec.Command("go", "build", "-o", c.Installer.PathToBin, "../observer/service_observer.go").
	// 	Output()
	// if err != nil {
	// 	return err
	// }

	// Creating config and installing daemons in system
	_, err = agentWork(*c)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot execute agent: %s", err.Error()))
	}

	_, err = observerWork(*c)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot execute observer: %s", err.Error()))
	}

	return nil
}

func observerWork(cfg InstallerConfig) (string, error) {
	observerOutput, err := exec.Command(
		cfg.Installer.PathToBin+"service_observer",
		srvDscFlg+cfg.Services.Observer.Description,
		srvNameFlg+cfg.Services.Observer.Name,
		srvDspNameFlg+cfg.Services.Observer.DisplayName,
		portFlg+cfg.Services.Observer.Port,
		"install",
	).Output()
	if err != nil {
		return "", err
	}
	return string(observerOutput), nil
}

func agentWork(cfg InstallerConfig) (string, error) {
	agentOutput, err := exec.Command(
		cfg.Installer.PathToBin+"service_observer",
		srvDscFlg+cfg.Services.Agent.Description,
		srvNameFlg+cfg.Services.Agent.Name,
		srvDspNameFlg+cfg.Services.Agent.DisplayName,
		portFlg+cfg.Services.Agent.Port,
		"install",
	).Output()
	if err != nil {
		return "", err
	}
	return string(agentOutput), nil
}

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
	cfg := new(InstallerConfig)

	err := cfg.Open(pathToCfg)
	if err != nil {
		log.Fatal("cannot open configuration file")
	}

	// Download files form server
	downloader.DownloadAllFiles(cfg.Files, cfg.ServerAddr+":"+cfg.ServerPort)

	// Start installation
	err = cfg.installation()
	if err != nil {
		log.Fatalf("Fatal error: %s\n", err.Error())
	}

	fmt.Println("Installer has finished work successfully")
}
