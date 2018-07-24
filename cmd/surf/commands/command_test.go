package commands_test

import (
	"github.com/CloudHub360/ch360.go/cmd/surf/commands"
	"github.com/CloudHub360/ch360.go/config"
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
		cmd:                new(commands.ListClassifiers).Usage(),
		expectedRunnerType: reflect.TypeOf(&commands.ListClassifiers{}),
	},
	{
		cmd:                new(commands.CreateClassifier).Usage(),
		expectedRunnerType: reflect.TypeOf(&commands.CreateClassifier{}),
	},
	{
		cmd:                new(commands.UploadClassifier).Usage(),
		expectedRunnerType: reflect.TypeOf(&commands.UploadClassifier{}),
	},
	{
		cmd:                new(commands.CreateExtractor).Usage(),
		expectedRunnerType: reflect.TypeOf(&commands.CreateExtractor{}),
	},
	{
		cmd:                new(commands.DeleteClassifier).Usage(),
		expectedRunnerType: reflect.TypeOf(&commands.DeleteClassifier{}),
	},
	{
		cmd:                new(commands.DeleteExtractor).Usage(),
		expectedRunnerType: reflect.TypeOf(&commands.DeleteExtractor{}),
	},
	{
		cmd:                new(commands.ListExtractors).Usage(),
		expectedRunnerType: reflect.TypeOf(&commands.ListExtractors{}),
	},
	{
		cmd:                new(commands.ClassifyCommand).Usage(),
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

		runParams := config.RunParams{}
		args.Bind(&runParams)

		// Act
		receivedRunner, _ := commands.CommandFor(&runParams)

		// Assert
		assert.Equal(t, testCase.expectedRunnerType, reflect.TypeOf(receivedRunner))
	}
}
