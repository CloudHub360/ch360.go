package types

type ClassificationResult struct {
	DocumentType       string
	IsConfident        bool
	RelativeConfidence float64
}
