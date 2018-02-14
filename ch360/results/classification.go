package results

type ClassificationResult struct {
	DocumentType       string
	IsConfident        bool
	RelativeConfidence float64
	DocumentTypeScores []DocumentTypeScore
}

type DocumentTypeScore struct {
	DocumentType string
	Score        float64
}
