package commands

import (
	"github.com/CloudHub360/ch360.go/config"
	"io"
)

type Login struct {
	configurationDirectory io.Writer
}

func NewLogin(configDirectory io.Writer) *Login {
	return &Login{
		configurationDirectory: configDirectory,
	}
}

func (cmd *Login) Execute(clientId string, clientSecret string) error {
	configuration := config.NewConfiguration(clientId, clientSecret)

	return configuration.Save(cmd.configurationDirectory)
}
