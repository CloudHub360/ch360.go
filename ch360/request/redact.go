package request

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
