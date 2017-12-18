package commands

import "github.com/howeyc/gopass"

type ConsoleSecretReader struct{}

func (reader *ConsoleSecretReader) Read() (string, error) {
	secret, err := gopass.GetPasswd()
	if err != nil {
		return "", err
	}

	return string(secret), nil
}
