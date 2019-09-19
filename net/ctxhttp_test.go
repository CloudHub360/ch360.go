package net_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/waives/surf/net"
	"github.com/waives/surf/net/mocks"
	"net/http"
	"testing"
)

func TestContextAwareHttpClient_Do_Calls_Underlying_And_Returns_Its_Values(t *testing.T) {
	// Arrange
	expectedResp := &http.Response{}
	expectedErr := errors.New("simulated error")
	mockHttpDoer := mocks.HttpDoer{}
	mockHttpDoer.On("Do", mock.Anything).Return(expectedResp, expectedErr)
	sut := net.NewContextAwareHttpClient(&mockHttpDoer)
	req := &http.Request{}

	// Act
	receivedResp, receivedErr := sut.Do(req)

	// Assert
	mockHttpDoer.AssertCalled(t, "Do", req)
	assert.Equal(t, expectedErr, receivedErr)
	assert.Equal(t, expectedResp, receivedResp)
}

func TestContextAwareHttpClient_Do_Returns_Err_From_Context_When_Cancelled(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	reqWithContext := (&http.Request{}).WithContext(ctx)
	fakeHttpDoer := FakeHttpDoer{}
	sut := net.NewContextAwareHttpClient(&fakeHttpDoer)

	var (
		receivedErr error
		expectedErr = context.Canceled
		reqDone     = make(chan bool)
	)

	// Act
	go func() {
		_, receivedErr = sut.Do(reqWithContext)
		reqDone <- true
	}()
	cancel()
	<-reqDone

	// Assert
	assert.Equal(t, expectedErr, receivedErr)
}

type FakeHttpDoer struct {
	resp *http.Response
	err  error
}

func (f *FakeHttpDoer) Do(req *http.Request) (*http.Response, error) {
	// Wait for context to be cancelled
	<-req.Context().Done()
	return f.resp, f.err
}
