package commands

type Login struct {
	client Getter
}

func NewLogin(client Getter) *Login {
	return &Login{
		client: client,
	}
}

func (cmd *Login) Execute(clientId string, clientSecret string) error {
	config := NewConfiguration(clientId, clientSecret)
	return config.Save()
}
