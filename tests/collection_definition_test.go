package stragollum_test

import (
	"encoding/json"
	"reflect"
	"stragollum/pkg/stragollum"
	"testing"
)

func TestCollectionDefinition_JSON(t *testing.T) {
	cd := stragollum.NewCollectionDefinition().
		WithDefaultID("string").
		WithIndexing(map[string]any{"foo": "bar"}).
		WithLexical("standard").
		WithRerank(&stragollum.RerankServiceOptions{
			Authentication: map[string]any{"k": "v"},
			ModelName:      "rerank-model",
			Parameters:     map[string]any{"k": "v"},
			Provider:       "rerank-provider",
		}).
		WithVector(&stragollum.CollectionVectorOptions{
			Dimension: ptrInt(999),
			Metric:    "cosine",
			Service: &stragollum.VectorServiceOptions{
				Authentication: map[string]any{"k": "v"},
				ModelName:      "vector-model",
				Parameters:     map[string]any{"k": "v"},
				Provider:       "vector-provider",
			},
			SourceModel: "source-model",
		})

	// Marshal to JSON
	actual, err := json.Marshal(cd)
	if err != nil {
		t.Fatalf("Failed to marshal CollectionDefinition: %v", err)
	}

	expectedJSON := `{"defaultId":{"type":"string"},"indexing":{"foo":"bar"},"lexical":{"analyzer":"standard","enabled":true},"rerank":{"enabled":true,"service":{"authentication":{"k":"v"},"modelName":"rerank-model","parameters":{"k":"v"},"provider":"rerank-provider"}},"vector":{"dimension":999,"metric":"cosine","service":{"authentication":{"k":"v"},"modelName":"vector-model","parameters":{"k":"v"},"provider":"vector-provider"},"sourceModel":"source-model"}}`

	var expected, got map[string]any
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}
	if err := json.Unmarshal(actual, &got); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("CollectionDefinition JSON does not match expected.\nExpected: %v\nGot: %v", expected, got)
	}
}

func ptrInt(i int) *int { return &i }
