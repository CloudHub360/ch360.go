package commands

import (
	"github.com/CloudHub360/ch360.go/ch360"
	"github.com/CloudHub360/ch360.go/cmd/ch360/commands/mocks"
	"github.com/stretchr/testify/assert"
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

	assert.Len(t, classifiersClient.Calls, 2)
	call := classifiersClient.Calls[1]

	assert.Len(t, call.Arguments, 1)
	actual := call.Arguments[0]

	assert.Equal(t, "charlie", actual)
}

func TestDeleteClassifier_Execute_Does_Not_Delete_The_Named_Classifier_When_It_Does_Not_Exist(t *testing.T) {
	classifiersClient := new(mocks.DeleterGetter)
	classifiersClient.On("GetAll", mock.Anything).Return(
		AListOfClassifiers("charlie", "jo", "chris"), nil)

	classifiersClient.On("Delete", mock.Anything).Return(nil)

	sut := NewDeleteClassifier(classifiersClient)

	sut.Execute("sydney")

	assert.Len(t, classifiersClient.Calls, 1)
	assert.Equal(t, classifiersClient.Calls[0].Method, "GetAll")
}

func AListOfClassifiers(names ...string) interface{} {
	expected := make(ch360.ClassifierList, len(names))

	for index, name := range names {
		expected[index] = ch360.Classifier{name}
	}

	return expected
}
