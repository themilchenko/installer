package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"

	"installer/downloader"
	"installer/utils"
)

const (
	defaultCfgContent = `services:
 agent:
   name: %s
   display_name: %s
   description: %s
   port: %s

 observer:
   name: %s
   display_name: %s
   description: %s
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
			Port        string `yaml:"port"`
		} `yaml:"agent"`

		Observer struct {
			Name        string `yaml:"name"`
			DisplayName string `yaml:"display_name"`
			Description string `yaml:"description"`
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

func (c *InstallerConfig) GenerateConfig() error {
	// Init agent service

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("========== Agent ==========")
	fmt.Print("Enter the Name of agent service: ")
	if scanner.Scan() {
		c.Services.Agent.Name = scanner.Text()
	}
	fmt.Print("Enter the Display Name of agent service: ")
	if scanner.Scan() {
		c.Services.Agent.DisplayName = scanner.Text()
	}
	fmt.Print("Enter the Description of agent service: ")
	if scanner.Scan() {
		c.Services.Agent.Description = scanner.Text()
	}
	fmt.Print("Enter the Port of agent service: ")
	fmt.Scanln(&c.Services.Agent.Port)

	fmt.Println()

	// Init observer service
	fmt.Println("========== Observer ==========")
	fmt.Print("Enter the Name of observer service: ")
	if scanner.Scan() {
		c.Services.Observer.Name = scanner.Text()
	}
	fmt.Print("Enter the Display Name of observer service: ")
	if scanner.Scan() {
		c.Services.Observer.DisplayName = scanner.Text()
	}
	fmt.Print("Enter the Description of observer service: ")
	if scanner.Scan() {
		c.Services.Observer.Description = scanner.Text()
	}
	fmt.Print("Enter the Port of observer service: ")
	fmt.Scanln(&c.Services.Observer.Port)

	fmt.Println()

	// Init pathes for installing services
	fmt.Println("========== Installation pathes ==========")
	fmt.Print("Enter the path to installing services configs: ")
	fmt.Scanln(&c.Installer.PathToCfg)
	fmt.Print("Enter the path to installing services binaries: ")
	fmt.Scanln(&c.Installer.PathToBin)

	fmt.Println()

	// Init name of files
	fmt.Println("========== Names to install ==========")
	fmt.Print("Enter the num of files you want to install from remote server: ")
	var filesNum int
	fmt.Scanln(&filesNum)
	fmt.Println("Name files you want to install")
	for i := 0; i < filesNum; i++ {
		var curName string
		fmt.Printf("%d) ", i+1)
		fmt.Scanln(&curName)
		c.Files = append(c.Files, curName)
	}

	fmt.Println()

	// Init the addr of server to install content from
	fmt.Println("========== Remote address ==========")
	fmt.Print("Enter the ip of remote server: ")
	fmt.Scanln(&c.ServerAddr)
	fmt.Print("Enter the port of remote server: ")
	fmt.Scanln(&c.ServerPort)

	// Marshal existing content
	src, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Create config
	dst, err := os.Create("config.yaml")
	if err != nil {
		return err
	}
	defer dst.Close()

	// Writing content in it
	buf := bytes.NewBuffer(src)

	_, err = io.Copy(dst, buf)
	if err != nil {
		return err
	}

	return nil
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
func (c *InstallerConfig) Installation() error {
	c.Installer.PathToCfg = getCorrectPath(c.Installer.PathToCfg)
	c.Installer.PathToBin = getCorrectPath(c.Installer.PathToBin)

	// Creating folder and config
	if _, err := os.Stat(c.Installer.PathToCfg + "aktiv.base/"); err != nil {
		if err := os.Mkdir(c.Installer.PathToCfg+"aktiv.base/", os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.Create(c.Installer.PathToCfg + "aktiv.base/config.yaml")
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
		c.Services.Agent.Port,

		c.Services.Observer.Name,
		c.Services.Observer.DisplayName,
		c.Services.Observer.Description,
		c.Services.Observer.Port,
	)
	_, err = fmt.Fprint(file, content)
	if err != nil {
		return err
	}

	// Copy binaries into global dir
	for _, file := range c.Files {
		if err = utils.CopyFile(downloader.DefaultDownloadPath+file, c.Installer.PathToBin+file); err != nil {
			return err
		}
		_, err := exec.Command("/usr/bin/chmod", "744", c.Installer.PathToBin+file).Output()
		if err != nil {
			return err
		}
	}

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
		"observer",
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
		"agent",
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
