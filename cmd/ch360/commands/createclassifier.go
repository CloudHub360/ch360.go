package commands

import "github.com/CloudHub360/ch360.go/ch360"

type CreateClassifier struct {
	client *ch360.ClassifiersClient
}

func NewCreateClassifier(client *ch360.ClassifiersClient) *CreateClassifier {
	return &CreateClassifier{
		client:client,
	}
}

func (cmd *CreateClassifier) Execute(classifierName string) error {
	return cmd.client.Create(classifierName)
}