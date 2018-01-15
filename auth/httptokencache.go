package auth

import (
	"sync"
	"time"
)

type TokenCache struct {
	retriever TokenRetriever
	token     *AccessToken
	once      sync.Once
	reqChan   chan bool
	respChan  chan tokenAndErr
}

type tokenAndErr struct {
	token *AccessToken
	err   error
}

func (cache *TokenCache) monitorRequestsForToken(reqChan chan bool, respChan chan tokenAndErr) {
	for range reqChan {
		if tokenIsFresh(cache.token) {
			respChan <- tokenAndErr{
				token: cache.token,
				err:   nil,
			}
			continue
		}

		token, err := cache.retriever.RetrieveToken()

		if err == nil {
			cache.token = token
		}

		respChan <- tokenAndErr{
			token: token,
			err:   err,
		}
	}
}

func (cache *TokenCache) RetrieveToken() (*AccessToken, error) {
	cache.once.Do(func() {
		go cache.monitorRequestsForToken(cache.reqChan, cache.respChan)
	})

	// Make a request to the monitoring goroutine to get a new token
	cache.reqChan <- true

	// Wait for a response
	res := <-cache.respChan

	return res.token, res.err
}

func tokenIsFresh(token *AccessToken) bool {
	if token == nil {
		return false
	}

	var expired = false
	deadline := time.Now().Add(time.Minute)
	expired = deadline.After(token.ExpiresAt)

	return !expired
}

func NewHttpTokenCache(tokenRetriever TokenRetriever) *TokenCache {

	return &TokenCache{
		retriever: tokenRetriever,
		reqChan:   make(chan bool),
		respChan:  make(chan tokenAndErr),
	}
}
