package commands

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/config"
	"io"
)

type Login struct {
	appDirectory   config.ConfigurationWriter
	tokenRetriever auth.TokenRetriever
	writer         io.Writer
}

func NewLogin(writer io.Writer, appDirectory config.ConfigurationWriter, retriever auth.TokenRetriever) *Login {
	return &Login{
		appDirectory:   appDirectory,
		tokenRetriever: retriever,
		writer:         writer,
	}
}

func (cmd *Login) Execute(clientId string, clientSecret string) error {

	fmt.Fprint(cmd.writer, "Logging in... ")

	err := cmd.execute(clientId, clientSecret)

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		fmt.Fprintln(cmd.writer, err.Error())
	} else {
		fmt.Fprintln(cmd.writer, "[OK]")
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
