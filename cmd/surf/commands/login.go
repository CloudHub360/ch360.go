package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

// ConfigureLoginCommand configures kingpin to add the login command.
func ConfigureLoginCommand(ctx context.Context,
	app *kingpin.Application,
	globalFlags *config.GlobalFlags) {
	app.Command("login", "Connect surf to your account.").
		Action(func(parseContext *kingpin.ParseContext) error {
			// execute the command
			return ExecuteWithMessage("Logging in... ", func() error {
				return execute(ctx, globalFlags)
			})
		})
}

func execute(ctx context.Context, flags *config.GlobalFlags) error {
	var (
		err          error
		clientId     = flags.ClientId
		clientSecret = flags.ClientSecret

		tokenRetriever = ch360.NewTokenRetriever(DefaultHttpClient, ch360.ApiAddress)
	)

	appDirectory, err := config.NewAppDirectory()
	if err != nil {
		return err
	}

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

	_, err = tokenRetriever.RetrieveToken(clientId, clientSecret)

	if err != nil {
		return err
	}

	configuration := config.NewConfiguration(clientId, clientSecret)

	err = appDirectory.WriteConfiguration(configuration)

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
