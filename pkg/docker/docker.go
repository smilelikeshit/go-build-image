package docker

import (
	"context"
	"docker-api/config"
	"docker-api/pkg/vault"
	"encoding/json"
	"fmt"
	"log"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	vaultapi "github.com/hashicorp/vault/api"
)

type DockerAuth struct {
	Auth          bool
	Username      *string `json:"USERNAME"`
	Password      *string `json:"PASSWORD"`
	DockerVersion *string
}

type HostConfig struct {
	TCP      string
	HostIP   string
	HostPort string
}

type Docker interface {
	Run(name, imagePath string, envVars []string, host HostConfig) (Container, error)
	ImageList() ([]docker.APIImages, error)
	BuildImage(name string, path string) error
	PushImage(name string) error
	DeleteImage(name string) error
	TagImage(name string) error
	Pull(name string) error
	HasImage(name string) (bool, error)
}

type dockerRepository struct {
	auth  DockerAuth
	vault *vaultapi.Client
	cli   *docker.Client
	ctx   context.Context
}

var newDockerClient = func() (*docker.Client, error) {
	return docker.NewClientFromEnv()
}

func NewDocker(auth *DockerAuth, config *config.Config) (Docker, error) {
	AuthDocker := DockerAuth{}
	var clientVault *vaultapi.Client

	cli, err := newDockerClient()
	if err != nil {
		return nil, err
	}

	docker := &dockerRepository{
		auth:  AuthDocker,
		vault: clientVault,
		cli:   cli,
		ctx:   context.Background(),
	}

	if auth.Auth == true {
		authVault := vault.VaultAuth{Username: &config.VaultUsername, Password: &config.VaultUsername}
		clientVault, err := vault.NewVault(config.VaultDNS, authVault)
		if err != nil {
			log.Println("Unable to create vault client: %s", err)
		}
		docker.vault = clientVault
		secret, err := clientVault.Logical().Read(config.VaultPath)
		if err != nil {
			fmt.Println(err)
		}

		// convert map to json
		jsonString, _ := json.Marshal(secret.Data)

		// convert json to struct
		json.Unmarshal(jsonString, &AuthDocker)
		docker.auth = AuthDocker
	}

	return docker, nil
}

func (d *dockerRepository) ImageList() ([]docker.APIImages, error) {
	imgs, err := d.cli.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return nil, err
	}

	return imgs, nil
}

func (d *dockerRepository) BuildImage(name string, path string) error {

	opts := docker.BuildImageOptions{
		Context:       d.ctx,
		Name:          name,
		Dockerfile:    "Dockerfile",
		ContextDir:    path,
		OutputStream:  os.Stdout,
		RawJSONStream: true,
	}

	if err := d.cli.BuildImage(opts); err != nil {
		return err
	}

	return nil
}

func (d *dockerRepository) PushImage(name string) error {

	repository, tag := docker.ParseRepositoryTag(name)
	opts := docker.PushImageOptions{
		Name:         repository,
		Tag:          tag,
		OutputStream: os.Stdout,
	}

	authconfiguration := docker.AuthConfiguration{
		Username: *d.auth.Username,
		Password: *d.auth.Password,
		// ServerAddress: *d.auth.DockerVersion,
	}

	if err := d.cli.PushImage(opts, authconfiguration); err != nil {
		return err
	}

	return nil
}

func (d *dockerRepository) DeleteImage(name string) error {
	if err := d.cli.RemoveImage(name); err != nil {
		return err
	}

	return nil
}

func (d *dockerRepository) TagImage(name string) error {
	repo, tag := docker.ParseRepositoryTag(name)
	config := docker.TagImageOptions{
		Repo:    repo,
		Tag:     tag,
		Force:   true,
		Context: d.ctx,
	}

	if err := d.cli.TagImage(name, config); err != nil {
		return err
	}

	return nil
}

func (d *dockerRepository) Pull(name string) error {
	config := docker.PullImageOptions{
		Registry:   "",
		Repository: name,
		Tag:        "latest",
	}
	if err := d.cli.PullImage(config, docker.AuthConfiguration{}); err != nil {
		return err
	}

	return nil
}

func (d *dockerRepository) HasImage(name string) (bool, error) {
	filters := map[string][]string{
		"reference": {name},
	}
	images, err := d.cli.ListImages(docker.ListImagesOptions{Filters: filters, Context: d.ctx})
	if err != nil {
		return false, err
	}
	return len(images) > 0, nil
}

func (d *dockerRepository) Run(name, imagePath string, envVars []string, host HostConfig) (Container, error) {
	image, err := d.HasImage(imagePath)
	if err != nil {
		return nil, err
	}

	if !image {
		if err := d.Pull(imagePath); err != nil {
			return nil, err
		}
	}

	return d.StartContainer(imagePath, name, envVars, host)
}

func (d *dockerRepository) StartContainer(imagePath, name string, envVars []string, host HostConfig) (Container, error) {

	portBindingss := map[docker.Port][]docker.PortBinding{
		docker.Port(host.TCP): {{HostIP: host.HostIP, HostPort: host.HostPort}}}

	portSets := map[docker.Port]struct{}{
		docker.Port(host.TCP): {}}

	createContHostConfig := docker.HostConfig{
		PortBindings:    portBindingss,
		PublishAllPorts: true,
		Privileged:      false,
	}

	createContConf := docker.Config{
		ExposedPorts: portSets,
		Image:        imagePath,
		Env:          envVars,
	}

	optsContainer := docker.CreateContainerOptions{
		Context:    d.ctx,
		Name:       name,
		Config:     &createContConf,
		HostConfig: &createContHostConfig,
	}

	container, err := d.cli.CreateContainer(optsContainer)
	if err != nil {
		log.Fatal(err)
	}
	if err := d.cli.StartContainerWithContext(container.ID, nil, d.ctx); err != nil {
		log.Fatal(err)
	}

	return newContainer(d.cli, d.ctx, container.ID), err
}
