package tests

import (
	"encoding/json"
	"github.com/CloudHub360/ch360.go/ch360/results"
)

func anExtractionResult() *results.ExtractionResult {
	var extractionResult results.ExtractionResult

	json.Unmarshal([]byte(anExtractionResponseJson), &extractionResult)

	return &extractionResult
}

var anExtractionResponseJsonList = "[" + anExtractionResponseJson + "]"

var anExtractionResponseJson = `{
  "field_results": [
    {
      "field_name": "Amount",
      "rejected": false,
      "reject_reason": "None",
      "result": {
        "text": "$5.50",
        "value": null,
        "rejected": false,
        "reject_reason": "None",
        "proximity_score": 100.0,
        "match_score": 100.0,
        "text_score": 100.0,
        "areas": [
          {
            "top": 558.7115,
            "left": 276.48,
            "bottom": 571.1989,
            "right": 298.58,
            "page_number": 1
          }
        ]
      },
      "alternative_results": null,
      "tabular_results": null
    }
  ],
  "page_sizes": {
    "page_count": 1,
    "pages": [
      {
        "page_number": 1,
        "width": 611.0,
        "height": 1008.0
      }
    ]
  }
}`
