package commands_test

import (
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/docopt/docopt-go"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

var expectedRunnerTypes = []struct {
	cmd                string
	expectedRunnerType reflect.Type
}{
	{
		cmd:                commands.ListClassifiersCommand,
		expectedRunnerType: reflect.TypeOf(&commands.ListClassifiers{}),
	},
	{
		cmd:                commands.CreateClassifierCommand,
		expectedRunnerType: reflect.TypeOf(&commands.CreateClassifier{}),
	},
	{
		cmd:                commands.CreateExtractorCommand,
		expectedRunnerType: reflect.TypeOf(&commands.CreateExtractor{}),
	},
	{
		cmd:                commands.DeleteClassifierCommand,
		expectedRunnerType: reflect.TypeOf(&commands.DeleteClassifier{}),
	},
	{
		cmd:                commands.ListExtractorsCommand,
		expectedRunnerType: reflect.TypeOf(&commands.ListExtractors{}),
	},
	{
		cmd:                commands.ClassifyFilesCommand,
		expectedRunnerType: reflect.TypeOf(&commands.ClassifyCommand{}),
	},
}

func TestRunnerBuilder_RunnerFor_Returns_Correct_Type(t *testing.T) {
	for _, testCase := range expectedRunnerTypes {
		// Arrange
		args := make(docopt.Opts)
		fields := strings.Fields(testCase.cmd)
		verb := fields[0]

		args[verb] = true

		if len(fields) > 1 {
			noun := fields[1]
			args[noun] = true
		}

		// Act
		receivedRunner, _ := commands.CommandFor(args)

		// Assert
		assert.Equal(t, testCase.expectedRunnerType, reflect.TypeOf(receivedRunner))
	}
}
