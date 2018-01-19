package tests

import (
	"bytes"
	"fmt"
	"github.com/CloudHub360/ch360.go/output/sinks"
	"github.com/CloudHub360/ch360.go/test/generators"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConsoleSinkSuite struct {
	suite.Suite
	sut      *sinks.BasicWriterSink
	contents string
	output   *bytes.Buffer
}

func (suite *ConsoleSinkSuite) SetupTest() {
	suite.output = &bytes.Buffer{}
	suite.contents = generators.String("contents")

	suite.sut = sinks.NewBasicWriterSink(suite.output)
}

func TestConsoleSinkRunner(t *testing.T) {
	suite.Run(t, new(ConsoleSinkSuite))
}

func (suite *ConsoleSinkSuite) TestOpen_Returns_Nil() {
	err := suite.sut.Open()

	assert.Nil(suite.T(), err)
}

func (suite *ConsoleSinkSuite) TestClose_Returns_Nil() {
	err := suite.sut.Open()

	assert.Nil(suite.T(), err)
}

func (suite *ConsoleSinkSuite) TestWrite_Delegates_To_Writer() {
	_, err := fmt.Fprint(suite.sut, suite.contents)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.contents, suite.output.String())
}

func (suite *ConsoleSinkSuite) TestWrite_Returns_Error_From_Writer() {
	sut := sinks.NewBasicWriterSink(&erroringWriter{})
	length, err := fmt.Fprint(sut, suite.contents)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), 0, length)
}

type erroringWriter struct{}

func (w *erroringWriter) Write(b []byte) (int, error) {
	return 0, errors.New("simulated error")
}
