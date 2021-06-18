package godata

import (
	"encoding/json"
	"testing"
)

type testResponseJson struct {
	ODataContext string              `json:"@odata.context"`
	ODataCount   int                 `json:"@odata.count"`
	Value        []testResponseValue `json:"value"`
}

type testResponseValue struct {
	Name   string
	Age    float64
	Gender string
}

func TestResponseWriter(t *testing.T) {

	test := &GoDataResponse{
		Fields: map[string]*GoDataResponseField{
			"@odata.context": {
				Value: "http://service.example",
			},
			"@odata.count": {
				Value: 8,
			},
			"value": {
				Value: []*GoDataResponseField{
					{
						Value: map[string]*GoDataResponseField{
							"Name": {Value: "John Doe"},
							"Age":  {11.01},
							"Male": {Value: "Female"},
						},
					},
					{
						Value: map[string]*GoDataResponseField{
							"Name":   {Value: "Jane \"Cool\" Doe"},
							"Age":    {12.1},
							"Gender": {Value: "Female"},
						},
					},
				},
			},
		},
	}

	written, err := test.Json()

	if err != nil {
		t.Error(err)
		return
	}

	var result testResponseJson
	err = json.Unmarshal(written, &result)

	if err != nil {
		t.Error(err)
		return
	}

	if result.ODataContext != "http://service.example" {
		t.Error("@odata.context is", result.ODataContext)
		return
	}

	if result.ODataCount != 8 {
		t.Error("@odata.count is", result.ODataCount)
		return
	}

	if len(result.Value) != 2 {
		t.Error("Result value is not length 2")
		return
	}

	if result.Value[0].Name != "John Doe" {
		t.Error("First value name is", result.Value[0].Name)
		return
	}

	if result.Value[1].Name != "Jane \"Cool\" Doe" {
		t.Error("Second value name is", result.Value[1].Name)
		return
	}

}
