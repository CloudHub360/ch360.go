package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type Login struct {
	client Getter
}

type configurationRoot struct {
	Configuration *configuration `json:"configuration"`
}

type configuration struct {
	Credentials apiCredentialsList `json:"credentials"`
}

type apiCredentialsList []apiCredentials

type apiCredentials struct {
	Context string `json:"context"`
	Url     string `json:"url"`
	Id      string `json:"client_id"`
	Secret  string `json:"client_secret"`
}

func NewLogin(client Getter) *Login {
	return &Login{
		client: client,
	}
}

func (cmd *Login) Execute(clientId string, clientSecret string) error {
	// Check credentials are valid (can get token)
	// TODO: Replace with a dedicated "can we get a token" call
	_, err := cmd.client.GetAll()

	if err != nil {
		// TODO: Better error message, or return err from dedicated check (see above)
		return errors.New("Invalid credentials")
	}

	// Store credentials to file
	var credentials = make(apiCredentialsList, 1)

	credentials[0] = apiCredentials{
		Id:     clientId,
		Secret: clientSecret,
	}

	config := &configuration{
		Credentials: credentials,
	}

	configRoot := &configurationRoot{
		Configuration: config,
	}

	json, _ := json.Marshal(configRoot)
	fmt.Println(string(json))

	CreateCH360DirIfNotExists()

	filename := filepath.Join(GetCH360Dir(), "config.json")

	err = ioutil.WriteFile(filename, json, 0644) //TODO: Permissions?

	return nil
}

func CreateCH360DirIfNotExists() {
	dir := GetCH360Dir()

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			os.Mkdir(dir, 0644) //TODO: Permissions?
		} else {
			// other error
			//TODO: Return error
		}
	}
}

func GetCH360Dir() string {
	return filepath.Join(UserHomeDir(), ".ch360")
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
