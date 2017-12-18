package commands

import (
	"fmt"
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
	classifiers, err := cmd.client.GetAll()
	if err != nil {
		fmt.Println("[FAILED]")
		return nil, err
	}

	return classifiers, nil
}
