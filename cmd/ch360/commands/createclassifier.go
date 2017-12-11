package commands

//go:generate mockery -name "Creator"

type Creator interface {
	Create(name string) error
}

type CreateClassifier struct {
	client Creator
}

func NewCreateClassifier(client Creator) *CreateClassifier {
	return &CreateClassifier{
		client: client,
	}
}

func (cmd *CreateClassifier) Execute(classifierName string) error {
	return cmd.client.Create(classifierName)
}
