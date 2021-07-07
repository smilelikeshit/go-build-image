package vault

import (
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

type VaultAuth struct {
	Username *string
	Password *string
	AppToken *string
}

func NewVault(dns string, auth VaultAuth) (*vaultapi.Client, error) {
	config := &vaultapi.Config{
		Address:    dns,
		Timeout:    time.Second * 60,
		MaxRetries: 2,
	}

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	if auth.AppToken != nil {
		client.SetToken(*auth.AppToken)
		return client, nil
	}
	token, err := UserPassLogin(client, *auth.Username, *auth.Password)
	if err != nil {
		return nil, err
	}
	client.SetToken(token)
	return client, nil
}

func UserPassLogin(client *vaultapi.Client, username string, password string) (string, error) {
	// to pass the password
	options := map[string]interface{}{
		"password": password,
	}
	path := fmt.Sprintf("auth/userpass/login/%s", username)

	// PUT call to get a token
	secret, err := client.Logical().Write(path, options)
	if err != nil {
		return "", err
	}
	token := secret.Auth.ClientToken
	return token, nil
}
