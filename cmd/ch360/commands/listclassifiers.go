package commands

import (
	"github.com/CloudHub360/ch360.go/ch360"
)

type ListClassifiers struct {
	client Getter
}

func NewListClassifiers(client Getter) *ListClassifiers {
	return &ListClassifiers{
		client: client,
	}
}

func (cmd *ListClassifiers) Execute() (ch360.ClassifierList, error) {
	return cmd.client.GetAll()
}
