package commands

import (
	"errors"
	"fmt"
	"github.com/CloudHub360/ch360.go/config"
	"github.com/docopt/docopt-go"
	"os"
)

type CredentialsResolver struct{}

func (r *CredentialsResolver) ResolveFromArgs(args docopt.Opts, configurationReader config.ConfigurationReader) (string, string, error) {
	var (
		clientIdArg, _     = args.String("--client-id")
		clientSecretArg, _ = args.String("--client-secret")
	)

	return r.Resolve(clientIdArg, clientSecretArg, configurationReader)
}

func (r *CredentialsResolver) Resolve(clientId string, clientSecret string, configurationReader config.ConfigurationReader) (string, string, error) {
	// If user specified both Id and Secret as parameters (or piping secret in), then use those values
	if clientId != "" && clientSecret != "" {
		return clientId, clientSecret, nil
	}

	// Specifying (or piping) just one of Id and Secret is not valid
	if (clientId != "" && clientSecret == "") || (clientId == "" && clientSecret != "") {
		return "", "", errors.New("You must either specify both --client-id and --client-secret, or neither.")
	}

	configuration, err := configurationReader.ReadConfiguration()
	if err != nil {
		if os.IsNotExist(err) {
			// Return sensible error if user hasn't logged in and there therefore is no
			// configuration file. This also masks other errors due to e.g. malformed
			// configuration file.
			return "", "", errors.New("Please run 'surf login' to connect to your account.")
		} else {
			return "", "", errors.New(fmt.Sprintf("There was an error loading your configuration file. Please run 'surf login' to connect to your account. Error: %s", err.Error()))
		}
	}

	if len(configuration.Credentials) == 0 {
		return "", "", errors.New("Your configuration file does not contain any credentials. Please run 'surf login' to connect to your account.")
	}

	clientId = configuration.Credentials[0].Id
	clientSecret = configuration.Credentials[0].Secret

	if clientId == "" || clientSecret == "" {
		return "", "", errors.New("Your configuration file does not contain valid credentials. Please run 'surf login' to connect to your account.")
	}

	return clientId, clientSecret, nil
}
