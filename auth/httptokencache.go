package auth

type tokenCache struct {
	retriever TokenRetriever
}

var cachedToken string

func newHttpTokenCache(tokenRetriever TokenRetriever) TokenRetriever {
	return &tokenCache{
		retriever: tokenRetriever,
	}
}

func (cache *tokenCache) RetrieveToken() (string, error) {
	var err error = nil
	if cachedToken == "" {
		cachedToken, err = cache.retriever.RetrieveToken()
	}

	return cachedToken, err
}
