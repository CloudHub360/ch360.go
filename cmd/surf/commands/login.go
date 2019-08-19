package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

// ConfigureLoginCommand configures kingpin to add the login command.
func ConfigureLoginCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {

	loginCmd := &LoginCmd{}
	app.Command("login", "Connect surf to your account.").
		Action(func(parseContext *kingpin.ParseContext) error {
			// execute the command
			return ExecuteWithMessage("Logging in... ", func() error {
				exitOnErr(loginCmd.initFromArgs(globalFlags))
				exitOnErr(loginCmd.Execute(ctx, globalFlags))

				return nil
			})
		})
}

type LoginCmd struct {
	TokenRetriever      auth.TokenRetriever
	ConfigurationWriter config.ConfigurationWriter
}

func (cmd *LoginCmd) initFromArgs(flags *config.GlobalFlags) error {

	var err error
	cmd.TokenRetriever = ch360.NewTokenRetriever(DefaultHttpClient, ch360.ApiAddress)
	cmd.ConfigurationWriter, err = config.NewAppDirectory()

	return err
}

func (cmd *LoginCmd) Execute(ctx context.Context, flags *config.GlobalFlags) error {
	var (
		err          error
		clientId     = flags.ClientId
		clientSecret = flags.ClientSecret
	)

	if clientId == "" {
		fmt.Print("\nClient Id: ")
		clientId, err = readSecretFromConsole()
		if err != nil {
			return err
		}
	}

	if clientSecret == "" {
		fmt.Print("Client Secret: ")
		clientSecret, err = readSecretFromConsole()
		if err != nil {
			return err
		}
	}

	_, err = cmd.TokenRetriever.RetrieveToken(clientId, clientSecret)

	if err != nil {
		return err
	}

	configuration := config.NewConfiguration(clientId, clientSecret)

	err = cmd.ConfigurationWriter.WriteConfiguration(configuration)

	return err
}

func readSecretFromConsole() (string, error) {
	secret, err := ConsoleSecretReader{}.Read()
	if err != nil {
		if err != ConsoleSecretReaderErrCancelled {
			fmt.Println(err.Error())
		}
		return "", err
	}
	return secret, nil
}
