package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/config"
	"io"
)

const LoginCommand = "login"

type Login struct {
	appDirectory   config.ConfigurationWriter
	tokenRetriever auth.TokenRetriever
	writer         io.Writer
	clientId       string
	clientSecret   string
}

func NewLoginFrom(runParams *config.RunParams, out io.Writer, appDir config.ConfigurationWriter, tokenRetriever auth.TokenRetriever) *Login {
	return NewLogin(out, appDir, tokenRetriever, runParams.ClientId, runParams.ClientSecret)
}

func NewLogin(out io.Writer, appDir config.ConfigurationWriter, tokenRetriever auth.TokenRetriever, clientId string, clientSecret string) *Login {
	return &Login{
		writer:         out,
		appDirectory:   appDir,
		clientId:       clientId,
		clientSecret:   clientSecret,
		tokenRetriever: tokenRetriever,
	}
}

func (cmd *Login) Execute(ctx context.Context) error {

	fmt.Fprint(cmd.writer, "Logging in... ")

	err := cmd.execute()

	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
	} else {
		fmt.Fprintln(cmd.writer, "[OK]")
	}
	return err
}

func (cmd *Login) execute() error {
	var (
		err error
	)

	if cmd.clientId == "" {
		fmt.Print("\nClient Id: ")
		cmd.clientId, err = cmd.readSecretFromConsole()
		if err != nil {
			return err
		}
	}

	if cmd.clientSecret == "" {
		fmt.Print("Client Secret: ")
		cmd.clientSecret, err = cmd.readSecretFromConsole()
		if err != nil {
			return err
		}
	}

	_, err = cmd.tokenRetriever.RetrieveToken(cmd.clientId, cmd.clientSecret)

	if err != nil {
		return err
	}

	configuration := config.NewConfiguration(cmd.clientId, cmd.clientSecret)

	err = cmd.appDirectory.WriteConfiguration(configuration)

	return err
}

func (cmd *Login) readSecretFromConsole() (string, error) {
	secret, err := ConsoleSecretReader{}.Read()
	if err != nil {
		if err != ConsoleSecretReaderErrCancelled {
			fmt.Println(err.Error())
		}
		return "", err
	}
	return secret, nil
}

func (cmd Login) Usage() string {
	return LoginCommand
}
