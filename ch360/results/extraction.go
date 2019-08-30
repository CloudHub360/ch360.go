package results

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
type ExtractForRedactionResult RedactedPdfRequest
type RedactedPdfRequest struct {
	Marks      []RedactionMark     `json:"marks"`
	ApplyMarks bool                `json:"apply_marks"`
	Bookmarks  []RedactionBookmark `json:"bookmarks"`
}
type RedactionArea struct {
	Top        float32 `json:"top"`
	Left       float32 `json:"left"`
	Bottom     float32 `json:"bottom"`
	Right      float32 `json:"right"`
	PageNumber float32 `json:"page_number"`
}
type RedactionMark struct {
	Area RedactionArea `json:"area"`
}
type RedactionBookmark struct {
	Text       string `json:"text"`
	PageNumber int    `json:"page_number"`
}
