package docker

import (
	"context"
	"docker-api/config"
	"docker-api/pkg/vault"
	"fmt"
	"log"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

var DefaultPushRetryCount = 3

type DockerAuth struct {
	VaultAuth     bool
	Username      *string
	Password      *string
	DockerVersion *string
}

type Docker interface {
	ImageList() ([]docker.APIImages, error)
	BuildImage(name string, path string) error
	PushImage(name string) error
	DeleteImage(name string) error
	TagImage(name string) error
}

type dockerRepository struct {
	auth DockerAuth
	cli  *docker.Client
	ctx  context.Context
}

var newDockerClient = func() (*docker.Client, error) {
	return docker.NewClientFromEnv()
}

func NewDocker(auth *DockerAuth, config *config.Config) (Docker, error) {
	AuthDocker := DockerAuth{}
	if auth.VaultAuth == true {
		fmt.Println(config)
		auth := vault.VaultAuth{Username: &config.VaultUsername, Password: &config.VaultUsername}
		clientVault, err := vault.NewVault(config.VaultDNS, auth)
		if err != nil {
			log.Println("Unable to create vault client: %s", err)
		}
		secret, err := clientVault.Logical().Read(config.VaultPath)
		if err != nil {
			fmt.Println(err)
		}
		AuthDocker.Username = secret.Data["username"].(*string)
		AuthDocker.Password = secret.Data["password"].(*string)
		AuthDocker.DockerVersion = secret.Data["docker_registry"].(*string)
	}
	cli, err := newDockerClient()
	if err != nil {
		return nil, err
	}

	docker := &dockerRepository{
		auth: AuthDocker,
		cli:  cli,
		ctx:  context.Background(),
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
		Username:      *d.auth.Username,
		Password:      *d.auth.Password,
		ServerAddress: *d.auth.DockerVersion,
	}

	for retries := 0; retries <= DefaultPushRetryCount; retries++ {
		if err := d.cli.PushImage(opts, authconfiguration); err != nil {
			return err
		}
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
