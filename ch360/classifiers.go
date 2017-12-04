package ch360


type ClassifiersClient struct {
	sender Sender
}

func (client *ClassifiersClient) CreateClassifier(name string) error {
	_, err := client.sender.Send("POST", "/classifiers/"+name, nil)
	return err
}
