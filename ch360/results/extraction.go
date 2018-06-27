package results

type FieldResult struct {
	FieldName    string `json:"field_name"`
	Rejected     bool   `json:"rejected"`
	RejectReason string `json:"reject_reason"`
	Result       *struct {
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
	} `json:"result"`
	AlternativeResults interface{} `json:"alternative_results"`
	TabularResults     interface{} `json:"tabular_results"`
}

type PageSizes struct {
	PageCount int `json:"page_count"`
	Pages     []struct {
		PageNumber int     `json:"page_number"`
		Width      float64 `json:"width"`
		Height     float64 `json:"height"`
	} `json:"pages"`
}

type ExtractionResult struct {
	FieldResults []FieldResult `json:"field_results"`
	PageSizes    PageSizes     `json:"page_sizes"`
}
