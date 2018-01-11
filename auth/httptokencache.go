package auth

import (
	"time"
)

type TokenCache struct {
	retriever TokenRetriever
	token     *AccessToken
}

func (cache *TokenCache) RetrieveToken() (*AccessToken, error) {
	if tokenIsFresh(cache.token) {
		return cache.token, nil
	}

	token, err := cache.retriever.RetrieveToken()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	cache.token = token
	return cache.token, nil
}

func tokenIsFresh(token *AccessToken) bool {
	if token == nil {
		return false
	}

	var expired = false
	deadline := time.Now().Add(-1 * time.Minute)
	expired = deadline.Before(token.ExpiresAt)

	return !expired
}

func NewHttpTokenCache(tokenRetriever TokenRetriever) *TokenCache {
	return &TokenCache{
		retriever: tokenRetriever,
	}
}
