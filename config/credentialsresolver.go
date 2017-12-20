package config

import (
	"errors"
	"fmt"
)

type CredentialsResolver struct{}

func (resolver *CredentialsResolver) Resolve(clientId string, clientSecret string, configurationReader ConfigurationReader) (string, string, error) {
	if clientId != "" && clientSecret != "" {
		return clientId, clientSecret, nil
	}

	configuration, err := configurationReader.ReadConfiguration()
	if err != nil {
		_, noConfigurationFile := err.(*NoConfigurationFileError)
		if noConfigurationFile {
			// Return sensible error if user hasn't logged in and there therefore is no
			// configuration file. This also masks other errors due to e.g. malformed
			// configuration file.
			return "", "", errors.New("Please run 'ch360 login' to setup your account.")
		} else {
			return "", "", errors.New(fmt.Sprintf("There was an error loading your configuration file. Please run 'ch360 login' to setup your account. Error: %s", err.Error()))
		}
	}

	if len(configuration.Credentials) == 0 {
		return "", "", errors.New("Your configuration file does not contain any credentials. Please run 'ch360 login' to setup your account.")
	}

	if clientId == "" {
		clientId = configuration.Credentials[0].Id
		if clientId == "" {
			return "", "", errors.New("Your configuration file does not contain an API Client Id. Please run 'ch360 login' to setup your account.")
		}
	}
	if clientSecret == "" {
		clientSecret = configuration.Credentials[0].Secret
		if clientSecret == "" {
			return "", "", errors.New("Your configuration file does not contain an API Client Secret. Please run 'ch360 login' to setup your account.")
		}
	}

	return clientId, clientSecret, nil
}
