package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type tokenCache struct {
	retriever TokenRetriever
	token     *jwt.Token
}

func newHttpTokenCache(tokenRetriever TokenRetriever) TokenRetriever {
	return &tokenCache{
		retriever: tokenRetriever,
	}
}

func (cache *tokenCache) RetrieveToken() (string, error) {
	if tokenIsFresh(cache.token) {
		return cache.token.Raw, nil
	}

	tokenString, err := cache.retriever.RetrieveToken()
	if err != nil {
		return "", err
	}

	token, err := parseToken(tokenString)
	if err != nil {
		return "", err
	}

	cache.token = token
	return cache.token.Raw, nil
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
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		deadline := time.Now().Add(-1 * time.Minute).Unix()
		expiry := getExpiry(claims)
		expired = deadline > expiry
	}

	return !expired
}

func getExpiry(claims jwt.MapClaims) int64 {
	switch exp := claims["exp"].(type) {
	case float64:
		return int64(exp)
	case json.Number:
		v, _ := exp.Int64()
		return v
	case string:
		v, _ := strconv.Atoi(exp)
		return int64(v)
	}

	return 0
}
