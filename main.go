package main

import (
	"log"
	"io/ioutil"
	"context"
	"bytes"
	"os"
	"archive/tar"
	"time"
	"fmt"
	"github.com/joho/godotenv"

	docker "github.com/fsouza/go-dockerclient"
	vaultapi "github.com/hashicorp/vault/api"

)

var DefaultPushRetryCount = 3


type DockerClient interface {
	BuildImage(opts docker.BuildImageOptions) error
	TagImage(name string, opts docker.TagImageOptions) error
	RemoveImage(name string) error
}

func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func newConfig() *vaultapi.Config {

	vault_addr := GetEnvVariable("VAULT_URL")

	config := &vaultapi.Config{
		Address: vault_addr,
		Timeout:      time.Second * 60,
		MaxRetries:   2,
	}

	return config
}


func removeImage(client *docker.Client, name string) error {
	if err := client.RemoveImage(name); err != nil {
		return err
	}

	return nil
}

func pushImage(client *docker.Client, name, username, password, url string) error{

	repository, tag := docker.ParseRepositoryTag(name)
	opts := docker.PushImageOptions{
		Name: repository,
		Tag:  tag,
		OutputStream: os.Stdout,
	}

	authconfiguration := docker.AuthConfiguration{
		Username : username,
		Password : password,
		ServerAddress : url,
		
	}
	
	for retries := 0; retries <= DefaultPushRetryCount; retries++ {
		if err := client.PushImage(opts,authconfiguration); err != nil {
			return err
		}
	}
	return nil
}


func buildImage(client DockerClient ,name, dockerfile string) error {


	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	inbuf := new(bytes.Buffer)
	tw := tar.NewWriter(inbuf)
	defer tw.Close()

	dockerFileReader, err := os.Open(dockerfile)
	if err != nil {
		return err
	}

	// Read the actual Dockerfile 
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		return err
	}

	// Make a TAR header for the file
	tarHeader := &tar.Header{
		Name: dockerfile,
		Size: int64(len(readDockerFile)),
	}

	// Writes the header described for the TAR file
	err = tw.WriteHeader(tarHeader)
    if err != nil {
		return err
    }

	// Writes the dockerfile data to the TAR file
    _, err = tw.Write(readDockerFile)
    if err != nil {
		return err
    }

    dockerFileTarReader := bytes.NewReader(inbuf.Bytes())

	
	buildImageOptions := docker.BuildImageOptions{
		Context : ctx, 
		Name : name, 
		Dockerfile : dockerfile,
		InputStream:  dockerFileTarReader,
		OutputStream: os.Stdout,
		RawJSONStream: true,
	}

	if err := client.BuildImage(buildImageOptions); err != nil {
		return err
	}

	return nil

}

func tagImage(dockerclient DockerClient, image, name string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	repo, tag := docker.ParseRepositoryTag(name)

	return dockerclient.TagImage(image, docker.TagImageOptions{
		Repo : repo, 
		Tag : tag, 
		Force : true, 
		Context: ctx,
	})

}


func main(){

	docker, err := docker.NewClientFromEnv()
	if err != nil {
		log.Println("Unable to create docker client: %s", err)
	}

	token := GetEnvVariable("VAULT_TOKEN")

	client, err := vaultapi.NewClient(newConfig())
	if err != nil {
		fmt.Println(err)
	}

	client.SetToken(token)

	secret, err := client.Logical().Read(GetEnvVariable("VAULT_PATH"))
	if err != nil {
		fmt.Println(err)

	}
	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		fmt.Printf("%T %#v\n", secret.Data["data"], secret.Data["data"])

	}


	dockerfile := GetEnvVariable("DOCKERFILE")

	name := GetEnvVariable("VERSION")

	err = buildImage(docker,name,dockerfile)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Has successfuly build %s ", name)
	}

	err = pushImage(docker,name,m["username"].(string),m["password"].(string),m["registry_url"].(string))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Has Successfully push to registry ")
	}

	err = removeImage(docker,name)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Clean image successfully")
	}

}