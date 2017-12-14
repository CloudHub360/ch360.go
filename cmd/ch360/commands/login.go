package commands

import "github.com/CloudHub360/ch360.go/config"

type Login struct {
	client Getter
}

func NewLogin(client Getter) *Login {
	return &Login{
		client: client,
	}
}

func (cmd *Login) Execute(clientId string, clientSecret string) error {
	config := config.NewConfiguration(clientId, clientSecret)
	return config.Save()
}
