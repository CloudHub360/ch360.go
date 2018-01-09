package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type tokenCache struct {
	retriever TokenRetriever
}

var cachedToken *jwt.Token

func newHttpTokenCache(tokenRetriever TokenRetriever) TokenRetriever {
	return &tokenCache{
		retriever: tokenRetriever,
	}
}

func (cache *tokenCache) RetrieveToken() (string, error) {
	if tokenIsFresh(cachedToken) {
		return cachedToken.Raw, nil
	}

	tokenString, err := cache.retriever.RetrieveToken()
	if err != nil {
		return "", err
	}

	token, err := parseToken(tokenString)
	if err != nil {
		return "", err
	}

	cachedToken = token
	return cachedToken.Raw, nil
}

func parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		kid := token.Header["kid"]
		if kid != nil {
			return kid, nil
		}

		return []byte("secret"), nil
	})
}

func tokenIsFresh(token *jwt.Token) bool {
	if token == nil || !token.Valid {
		return false
	}

	var expired = false
	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		deadline := time.Now().Add(-1 * time.Minute)
		expired = claims.VerifyExpiresAt(int64(deadline.Unix()), true)
	}

	return !expired
}
