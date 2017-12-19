package commands

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/config"
	"os"
)

//go:generate mockery -name Reader
type Reader interface {
	Read() (string, error)
}

type Login struct {
	appDirectory config.ConfigurationWriter
	secretReader Reader
}

func NewLogin(appDirectory config.ConfigurationWriter, reader Reader) *Login {
	return &Login{
		appDirectory: appDirectory,
		secretReader: reader,
	}
}

func (cmd *Login) Execute(clientId string, clientSecret string) error {
	var err error
	if clientSecret == "" {
		clientSecret, err = cmd.readSecret()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	fmt.Print("Logging in... ")
	configuration := config.NewConfiguration(clientId, clientSecret)

	err = cmd.appDirectory.WriteConfiguration(configuration)
	if err != nil {
		fmt.Println("[FAILED]")
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		fmt.Println("[OK]")
	}
	return err
}

func (cmd *Login) readSecret() (string, error) {
	fmt.Print("API Client Secret: ")
	secret, err := cmd.secretReader.Read()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return secret, nil
}
