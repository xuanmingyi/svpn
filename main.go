package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "", log.LstdFlags)
}

type ServiceConfig struct {
	Type   string `yaml:"type"`
	Listen string `yaml:"listen"`
	Port   int    `yaml:"port"`
	Device string `yaml:"device"`
	Host   string `yaml:"host"`
}

func InitService(path string) (service *ServiceConfig, err error) {
	var content []byte
	service = &ServiceConfig{}

	content, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, service)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func main() {
	if len(os.Args) != 2 {
		logger.Println("hello world")
	}

	service, err := InitService(os.Args[1])
	if err != nil {
		logger.Println("service error")
	}

	if service.Type == "server" {
		server, _ := InitServer(service)

		server.Init()

		server.Run()
	} else if service.Type == "client" {
		client, _ := InitClient(service)

		client.Init()

		client.Run()
	}
}
