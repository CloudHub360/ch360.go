package results

import "github.com/CloudHub360/ch360.go/ch360/request"

type InnerResult struct {
	Text           string      `json:"text"`
	Value          interface{} `json:"value"`
	Rejected       bool        `json:"rejected"`
	RejectReason   string      `json:"reject_reason"`
	ProximityScore float64     `json:"proximity_score"`
	MatchScore     float64     `json:"match_score"`
	TextScore      float64     `json:"text_score"`
	Areas          []struct {
		Top        float64 `json:"top"`
		Left       float64 `json:"left"`
		Bottom     float64 `json:"bottom"`
		Right      float64 `json:"right"`
		PageNumber int     `json:"page_number"`
	} `json:"areas"`
}

type FieldResult struct {
	FieldName          string         `json:"field_name"`
	Rejected           bool           `json:"rejected"`
	RejectReason       string         `json:"reject_reason"`
	Result             *InnerResult   `json:"result"`
	AlternativeResults []*InnerResult `json:"alternatives"`
	TabularResults     interface{}    `json:"tabular_results"`
}

type Document struct {
	PageCount int `json:"page_count"`
	Pages     []struct {
		PageNumber int     `json:"page_number"`
		Width      float64 `json:"width"`
		Height     float64 `json:"height"`
	} `json:"pages"`
}

type ExtractionResult struct {
	FieldResults []FieldResult `json:"field_results"`
	Document     Document      `json:"document"`
}

// ExtractForRedactionResult represents the json returned by waives' 'extract document' endpoint
// when the 'Accept' header is set to 'application/vnd.waives.requestformats.redact+json' - ie,
// the same format as the 'get redact PDF' request.
type ExtractForRedactionResult request.RedactedPdfRequest
