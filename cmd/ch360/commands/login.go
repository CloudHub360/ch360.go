package commands

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/config"
	"os"
)

//go:generate mockery -name Reader
type Reader interface {
	Read() (string, error)
}

type Login struct {
	appDirectory   config.ConfigurationWriter
	secretReader   Reader
	tokenRetriever auth.TokenRetriever
}

func NewLogin(appDirectory config.ConfigurationWriter, reader Reader, retriever auth.TokenRetriever) *Login {
	return &Login{
		appDirectory:   appDirectory,
		secretReader:   reader,
		tokenRetriever: retriever,
	}
}

func (cmd *Login) Execute(clientId string, clientSecret string) error {

	fmt.Print("Logging in... ")

	err := cmd.execute(clientId, clientSecret)

	if err != nil {
		fmt.Println("[FAILED]")
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		fmt.Println("[OK]")
	}
	return err
}

func (cmd *Login) execute(clientId string, clientSecret string) error {
	var err error

	if clientSecret == "" {
		clientSecret, err = cmd.readSecret()
		if err != nil {
			return err
		}
	}

	_, err = cmd.tokenRetriever.RetrieveToken()

	if err != nil {
		return err
	}

	configuration := config.NewConfiguration(clientId, clientSecret)

	err = cmd.appDirectory.WriteConfiguration(configuration)

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
