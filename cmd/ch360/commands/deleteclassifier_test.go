package commands

import (
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeleteClassifier_Execute_Deletes_The_Named_Classifier_When_It_Exists(t *testing.T) {
	classifiersClient := new(mocks.DeleterGetter)
	classifiersClient.On("GetAll", mock.Anything).Return(
		AListOfClassifiers("charlie", "jo", "chris"), nil)

	classifiersClient.On("Delete", mock.Anything).Return(nil)

	sut := NewDeleteClassifier(classifiersClient)

	sut.Execute("charlie")

	classifiersClient.AssertCalled(t, "GetAll")
	classifiersClient.AssertCalled(t, "Delete", "charlie")
}

func TestDeleteClassifier_Execute_Does_Not_Delete_The_Named_Classifier_When_It_Does_Not_Exist(t *testing.T) {
	classifiersClient := new(mocks.DeleterGetter)
	classifiersClient.On("GetAll", mock.Anything).Return(
		AListOfClassifiers("charlie", "jo", "chris"), nil)

	sut := NewDeleteClassifier(classifiersClient)

	sut.Execute("sydney")

	classifiersClient.AssertCalled(t, "GetAll")
	classifiersClient.AssertNotCalled(t, "Delete")
}

func AListOfClassifiers(names ...string) interface{} {
	expected := make(ch360.ClassifierList, len(names))

	for index, name := range names {
		expected[index] = ch360.Classifier{name}
	}

	return expected
}
