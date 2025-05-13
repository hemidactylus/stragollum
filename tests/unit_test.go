package tqests

import (
	"stragollum/pkg/stragollum"
	"testing"
)

func TestAstraDBClient(t *testing.T) {
	client := stragollum.NewAstraDBClient("test_env", "test_tok")
	result := client.GetEnvironment()
	expected := "test_env"

	if result != expected {
		t.Errorf("GetEnvironment() = %s; want %s", result, expected)
	}
}
