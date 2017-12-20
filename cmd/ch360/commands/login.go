package commands

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/config"
)

type Login struct {
	appDirectory   config.ConfigurationWriter
	tokenRetriever auth.TokenRetriever
}

func NewLogin(appDirectory config.ConfigurationWriter, retriever auth.TokenRetriever) *Login {
	return &Login{
		appDirectory:   appDirectory,
		tokenRetriever: retriever,
	}
}

func (cmd *Login) Execute(clientId string, clientSecret string) error {

	fmt.Print("Logging in... ")

	err := cmd.execute(clientId, clientSecret)

	if err != nil {
		fmt.Println("[FAILED]")
		fmt.Println(err.Error())
	} else {
		fmt.Println("[OK]")
	}
	return err
}

func (cmd *Login) execute(clientId string, clientSecret string) error {
	var err error

	_, err = cmd.tokenRetriever.RetrieveToken()

	if err != nil {
		return err
	}

	configuration := config.NewConfiguration(clientId, clientSecret)

	err = cmd.appDirectory.WriteConfiguration(configuration)

	return err
}
