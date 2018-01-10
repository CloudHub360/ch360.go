package auth

import (
	"github.com/CloudHub360/ch360.go/auth/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

const EXPIRED_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImV4cCI6IjE1MTU0MTE5ODIifQ.efsGk6oZDo3PK5euKvuoa-KDHcXY5gQUGdoeN-OO9LA"

// Expires 9 Jan 2118
const VALID_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImV4cCI6IjQ2NzExNzE5ODIifQ.vsp4YkbzAwwogc9qCjwjICSXqVARjVKL6neEm7iHnYY"

func Test_RetrieveToken_Uses_Cached_Token_If_It_Has_Not_Expired(t *testing.T) {
	var tokenRetriever = new(mocks.TokenRetriever)
	var sut = newHttpTokenCache(tokenRetriever)

	tokenRetriever.On("RetrieveToken").Return(VALID_TOKEN, nil)
	sut.RetrieveToken()

	token, err := sut.RetrieveToken()

	assert.Nil(t, err)
	assert.Equal(t, VALID_TOKEN, token)
	tokenRetriever.AssertNumberOfCalls(t, "RetrieveToken", 1)
}

func Test_RetrieveToken_Requests_New_Token_If_It_Has_Expired(t *testing.T) {
	var tokenRetriever = new(mocks.TokenRetriever)
	var sut = newHttpTokenCache(tokenRetriever)

	tokenRetriever.On("RetrieveToken").Return(EXPIRED_TOKEN, nil)
	sut.RetrieveToken()
	tokenRetriever.ExpectedCalls = nil
	tokenRetriever.On("RetrieveToken").Return(VALID_TOKEN, nil)

	token, err := sut.RetrieveToken()

	assert.Nil(t, err)
	assert.Equal(t, VALID_TOKEN, token)
	tokenRetriever.AssertNumberOfCalls(t, "RetrieveToken", 2)
}
