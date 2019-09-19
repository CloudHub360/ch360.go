package net

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/waives/surf/ioutils"
	"github.com/waives/surf/net/mocks"
	"net/http"
	"testing"
)

func TestRetryingHttpClient_Should_Return_Response_From_Wrapped_Doer_On_Success(t *testing.T) {
	// Arrange
	var (
		wrappedDoer            = mocks.HttpDoer{}
		expectedResponse       = &http.Response{StatusCode: 200}
		expectedErr      error = nil
		request, _             = http.NewRequest("GET", "https://api.waives.io/version", nil)
	)
	wrappedDoer.On("Do", mock.Anything).Return(expectedResponse, expectedErr)
	sut := NewRetryingHttpClient(&wrappedDoer, 0, 0.01)

	// Act
	actualResponse, actualErr := sut.Do(request)

	// Assert
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Equal(t, expectedErr, actualErr)
}

func TestRetryingHttpClient_Should_Retry_On_Network_Error(t *testing.T) {
	// Arrange
	var (
		wrappedDoer                      = mocks.HttpDoer{}
		expectedResponse  *http.Response = nil
		expectedErr                      = errors.New("test error")
		actualCallCount   int
		retryAttempts     = 3
		expectedCallCount = retryAttempts + 1
		request, _        = http.NewRequest("GET", "https://api.waives.io/version", nil)
	)
	wrappedDoer.
		On("Do", mock.Anything).
		Run(func(_ mock.Arguments) {
			actualCallCount++
		}).
		Return(expectedResponse, expectedErr)
	sut := NewRetryingHttpClient(&wrappedDoer, retryAttempts, 0.01)

	// Act
	_, _ = sut.Do(request)

	// Assert
	assert.Equal(t, expectedCallCount, actualCallCount)
}

func TestRetryingHttpClient_Should_Retry_On_Http_500_And_Over(t *testing.T) {
	// Arrange
	var (
		wrappedDoer             = mocks.HttpDoer{}
		expectedResponse        = &http.Response{StatusCode: 500}
		expectedErr       error = nil
		retryAttempts           = 3
		expectedCallCount       = retryAttempts + 1
		request, _              = http.NewRequest("GET", "https://api.waives.io/version", nil)
		actualCallCount   int
	)
	wrappedDoer.
		On("Do", mock.Anything).
		Run(func(_ mock.Arguments) {
			actualCallCount++
		}).
		Return(expectedResponse, expectedErr)
	sut := NewRetryingHttpClient(&wrappedDoer, retryAttempts, 0.01)

	// Act
	_, _ = sut.Do(request)

	// Assert
	assert.Equal(t, expectedCallCount, actualCallCount)
}

func TestRetryingHttpClient_Should_Specify_The_Same_Request_Data_On_Each_Retry(t *testing.T) {
	// Arrange
	var (
		wrappedDoer                        = mocks.HttpDoer{}
		expectedResponse    *http.Response = nil
		expectedErr                        = errors.New("test error")
		expectedBody                       = []byte("test request body")
		retryAttempts                      = 1
		request, _                         = http.NewRequest("GET", "https://api.waives.io/version", bytes.NewBuffer(expectedBody))
		actualRequests      []*http.Request
		actualRequestBodies [][]byte
	)
	wrappedDoer.
		On("Do", mock.Anything).
		Run(func(args mock.Arguments) {
			// Capture request and request body on each attempt
			actualRequest := (args.Get(0)).(*http.Request)
			actualRequests = append(actualRequests, actualRequest)
			bodyBuf, _ := ioutils.DrainClose(request.Body)
			actualRequestBodies = append(actualRequestBodies, bodyBuf.Bytes())
		}).
		Return(expectedResponse, expectedErr)
	sut := NewRetryingHttpClient(&wrappedDoer, retryAttempts, 0.01)

	// Act
	_, _ = sut.Do(request)

	// Assert
	for _, actualRequest := range actualRequests {
		assert.Equal(t, request, actualRequest)
	}
	for _, actualRequestBody := range actualRequestBodies {
		assert.Equal(t, expectedBody, actualRequestBody)
	}
}

func TestRetryingHttpClient_Should_Not_Retry_If_Cancelled(t *testing.T) {
	// Arrange
	var (
		wrappedDoer                      = mocks.HttpDoer{}
		expectedResponse  *http.Response = nil
		retryAttempts                    = 1
		request, _                       = http.NewRequest("GET", "https://api.waives.io/version", nil)
		actualCallCount                  = 0
		expectedCallCount                = 1
		ctx, cancelFn                    = context.WithCancel(context.Background())
	)
	cancelFn() // cancel request's Context
	expectedErr := ctx.Err()
	wrappedDoer.
		On("Do", mock.Anything).
		Run(func(_ mock.Arguments) {
			actualCallCount++
		}).
		Return(expectedResponse, expectedErr)

	request = request.WithContext(ctx)
	sut := NewRetryingHttpClient(&wrappedDoer, retryAttempts, 0.01)

	// Act
	_, actualErr := sut.Do(request)

	// Assert
	assert.Equal(t, expectedCallCount, actualCallCount)
	assert.Equal(t, expectedErr, actualErr)
}
