package main

import (
	"docker-api/config"
	"docker-api/pkg/docker"
	"fmt"
	"log"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	client, err := docker.NewDocker(&docker.DockerAuth{VaultAuth: true}, &config)
	if err != nil {
		log.Println("Unable to create docker client: %s", err)
	}

	NameImage := "test"
	PathProject := "example"

	//build image docker
	if err = client.BuildImage(NameImage, PathProject); err != nil {
		fmt.Println(err)
	}

	//push image docker
	err = client.PushImage(NameImage)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Has Successfully push to registry ")
	}
}
