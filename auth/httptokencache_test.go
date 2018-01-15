package auth_test

import (
	"github.com/CloudHub360/ch360.go/auth"
	"github.com/CloudHub360/ch360.go/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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
}

func (suite *tokenCacheSuite) SetupTest() {
	suite.tokenRetriever = new(mocks.TokenRetriever)
	suite.sut = auth.NewHttpTokenCache(suite.tokenRetriever)
}

func TestTokenCacheSuiteRunner(t *testing.T) {
	suite.Run(t, new(tokenCacheSuite))
}

func (suite *tokenCacheSuite) Test_RetrieveToken_Uses_Cached_Token_If_It_Has_Not_Expired() {
	suite.populateCache(validToken)

	token, err := suite.sut.RetrieveToken()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), &validToken, token)
	suite.tokenRetriever.AssertNumberOfCalls(suite.T(), "RetrieveToken", 1)
}

func (suite *tokenCacheSuite) Test_RetrieveToken_Requests_New_Token_If_It_Has_Expired() {
	suite.populateCache(expiredToken)
	suite.reStub("RetrieveToken").Return(&validToken, nil)

	token, err := suite.sut.RetrieveToken()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), &validToken, token)
	suite.tokenRetriever.AssertNumberOfCalls(suite.T(), "RetrieveToken", 2)
}

func (suite *tokenCacheSuite) Test_RetrieveToken_Only_Requests_New_Token_Once_If_Used_In_Parallel() {
	suite.reStub("RetrieveToken").Return(&validToken, nil)

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			suite.sut.RetrieveToken()
			wg.Done()
		}()
	}

	wg.Wait()

	suite.tokenRetriever.AssertNumberOfCalls(suite.T(), "RetrieveToken", 1)
}

func (suite *tokenCacheSuite) populateCache(token auth.AccessToken) {
	suite.reStub("RetrieveToken").Return(&token, nil)
	suite.sut.RetrieveToken()
}

func (suite *tokenCacheSuite) reStub(methodName string) *mock.Call {
	suite.tokenRetriever.ExpectedCalls = nil
	return suite.tokenRetriever.On(methodName)
}
