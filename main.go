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

	client, err := docker.NewDocker(&docker.DockerAuth{Auth: true}, &config)
	if err != nil {
		log.Println("Unable to create docker client: %s", err)
	}

	NameImage := "arf95/test"
	PathProject := "example"

	// build image docker
	if err = client.BuildImage(NameImage, PathProject); err != nil {
		fmt.Println(err)
	}

	// tag image docker
	if err = client.TagImage(NameImage); err != nil {
		fmt.Println(err)
	}

	// push image docker
	err = client.PushImage(NameImage)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Has Successfully push to registry ")
	}

	// run container
	client.Run("test-container", NameImage, nil, docker.HostConfig{TCP: "8080/tcp", HostIP: "127.0.0.1", HostPort: "8080"})

	// delete
	if err := client.DeleteImage(NameImage); err != nil {
		fmt.Println(err)
	}
}
