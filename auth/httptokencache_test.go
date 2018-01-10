package auth

import (
	"github.com/CloudHub360/ch360.go/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

const EXPIRED_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImV4cCI6IjE1MTU0MTE5ODIifQ.efsGk6oZDo3PK5euKvuoa-KDHcXY5gQUGdoeN-OO9LA"

// Expires 9 Jan 2118
const VALID_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImV4cCI6IjQ2NzExNzE5ODIifQ.vsp4YkbzAwwogc9qCjwjICSXqVARjVKL6neEm7iHnYY"

type tokenCacheSuite struct {
	suite.Suite
	tokenRetriever *mocks.TokenRetriever
	sut            TokenRetriever
}

func (suite *tokenCacheSuite) SetupTest() {
	suite.tokenRetriever = new(mocks.TokenRetriever)
	suite.sut = newHttpTokenCache(suite.tokenRetriever)

	suite.tokenRetriever.On("RetrieveToken").Return(VALID_TOKEN, nil)
}

func TestTokenCacheSuiteRunner(t *testing.T) {
	suite.Run(t, new(tokenCacheSuite))
}

func (suite *tokenCacheSuite) Test_RetrieveToken_Uses_Cached_Token_If_It_Has_Not_Expired() {
	suite.sut.RetrieveToken()

	token, err := suite.sut.RetrieveToken()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), VALID_TOKEN, token)
	suite.tokenRetriever.AssertNumberOfCalls(suite.T(), "RetrieveToken", 1)
}

func (suite *tokenCacheSuite) Test_RetrieveToken_Requests_New_Token_If_It_Has_Expired() {
	suite.reStub("RetrieveToken").Return(EXPIRED_TOKEN, nil)
	suite.sut.RetrieveToken()
	suite.reStub("RetrieveToken").Return(VALID_TOKEN, nil)

	token, err := suite.sut.RetrieveToken()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), VALID_TOKEN, token)
	suite.tokenRetriever.AssertNumberOfCalls(suite.T(), "RetrieveToken", 2)
}

func (suite *tokenCacheSuite) reStub(methodName string, returnArguments ...interface{}) *mock.Call {
	suite.tokenRetriever.ExpectedCalls = nil
	return suite.tokenRetriever.On(methodName)
}
