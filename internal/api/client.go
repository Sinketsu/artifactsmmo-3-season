package api

import "github.com/Sinketsu/artifactsmmo-3-season/gen/oas"

type Client struct {
	*oas.Client
}

func New(serverUrl string, token string) (*Client, error) {
	auth := &Auth{Token: token}
	cli, err := oas.NewClient(serverUrl, auth)
	return &Client{Client: cli}, err
}
