package auth_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/waives/surf/auth"
	"github.com/waives/surf/auth/mocks"
	"github.com/waives/surf/test/generators"
	"sync"
	"testing"
	"time"
)

var expiredToken = auth.AccessToken{
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImV4cCI6IjE1MTU0MTE5ODIifQ.efsGk6oZDo3PK5euKvuoa-KDHcXY5gQUGdoeN-OO9LA",
	time.Date(2018, time.Month(1), 8, 11, 46, 22, 0, &time.Location{}),
}
var validToken = auth.AccessToken{
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImV4cCI6IjQ2NzExNzE5ODIifQ.vsp4YkbzAwwogc9qCjwjICSXqVARjVKL6neEm7iHnYY",
	time.Date(2118, time.Month(1), 9, 11, 46, 22, 0, &time.Location{}),
}

type tokenCacheSuite struct {
	suite.Suite
	tokenRetriever *mocks.TokenRetriever
	sut            auth.TokenRetriever
	clientId       string
	clientSecret   string
}

func (suite *tokenCacheSuite) SetupTest() {
	suite.tokenRetriever = new(mocks.TokenRetriever)
	suite.sut = auth.NewHttpTokenCache(suite.tokenRetriever)

	suite.clientId = generators.String("client-id")
	suite.clientSecret = generators.String("client-secret")
}

func TestTokenCacheSuiteRunner(t *testing.T) {
	suite.Run(t, new(tokenCacheSuite))
}

func (suite *tokenCacheSuite) Test_RetrieveToken_Uses_Cached_Token_If_It_Has_Not_Expired() {
	suite.populateCache(validToken)

	token, err := suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), &validToken, token)
	suite.tokenRetriever.AssertNumberOfCalls(suite.T(), "RetrieveToken", 1)
}

func (suite *tokenCacheSuite) Test_RetrieveToken_Requests_New_Token_If_It_Has_Expired() {
	suite.populateCache(expiredToken)
	suite.reStub("RetrieveToken", mock.Anything, mock.Anything).Return(&validToken, nil)

	token, err := suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), &validToken, token)
	suite.tokenRetriever.AssertNumberOfCalls(suite.T(), "RetrieveToken", 2)
}

func (suite *tokenCacheSuite) Test_RetrieveToken_Only_Requests_New_Token_Once_If_Used_In_Parallel() {
	suite.reStub("RetrieveToken", mock.Anything, mock.Anything).Return(&validToken, nil)

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)
			wg.Done()
		}()
	}

	wg.Wait()

	suite.tokenRetriever.AssertNumberOfCalls(suite.T(), "RetrieveToken", 1)
}

func (suite *tokenCacheSuite) populateCache(token auth.AccessToken) {
	suite.reStub("RetrieveToken", mock.Anything, mock.Anything).Return(&token, nil)
	suite.sut.RetrieveToken(suite.clientId, suite.clientSecret)
}

func (suite *tokenCacheSuite) reStub(methodName string, arguments ...interface{}) *mock.Call {
	suite.tokenRetriever.ExpectedCalls = nil
	return suite.tokenRetriever.On(methodName, arguments...)
}
