package tests

import (
	"github.com/CloudHub360/ch360.go/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RunParams_Validate_Returns_Err_If_Both_Single_And_Multi_File_Enabled(t *testing.T) {
	// Arrange
	sut := config.RunParams{
		OutputFile:   "singlefile",
		MultiFileOut: true,
	}

	// Act
	err := sut.Validate()

	// Assert
	assert.NotNil(t, err)
}

func Test_RunParams_Noun(t *testing.T) {
	var nounTests = []struct {
		Extractor    bool
		Extractors   bool
		Classifier   bool
		Classifiers  bool
		ExpectedNoun config.Noun
	}{
		{
			Extractor:    true,
			ExpectedNoun: config.Extractor,
		},
		{
			Extractors:   true,
			ExpectedNoun: config.Extractor,
		},
		{
			Classifier:   true,
			ExpectedNoun: config.Classifier,
		},
		{
			Classifiers:  true,
			ExpectedNoun: config.Classifier,
		},
		{
			ExpectedNoun: -1,
		},
	}

	for _, nounTest := range nounTests {
		sut := config.RunParams{}
		sut.Extractor = nounTest.Extractor
		sut.Extractors = nounTest.Extractors
		sut.Classifier = nounTest.Classifier
		sut.Classifiers = nounTest.Classifiers

		assert.Equal(t, nounTest.ExpectedNoun, sut.Noun())
	}
}

func Test_RunParams_Verb(t *testing.T) {
	var nounTests = []struct {
		Extractor    bool
		Extractors   bool
		Classifier   bool
		Classifiers  bool
		ExpectedNoun config.Noun
	}{
		{
			Extractor:    true,
			ExpectedNoun: config.Extractor,
		},
		{
			Extractors:   true,
			ExpectedNoun: config.Extractor,
		},
		{
			Classifier:   true,
			ExpectedNoun: config.Classifier,
		},
		{
			Classifiers:  true,
			ExpectedNoun: config.Classifier,
		},
		{
			ExpectedNoun: -1,
		},
	}

	for _, nounTest := range nounTests {
		sut := config.RunParams{}
		sut.Extractor = nounTest.Extractor
		sut.Extractors = nounTest.Extractors
		sut.Classifier = nounTest.Classifier
		sut.Classifiers = nounTest.Classifiers

		assert.Equal(t, nounTest.ExpectedNoun, sut.Noun())
	}
}
