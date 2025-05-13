package tqests

import (
	"stragollum/pkg/stragollum"
	"testing"
)

func TestAdd(t *testing.T) {
	result := stragollum.Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}

func TestAstraDBClient(t *testing.T) {
	client := stragollum.NewAstraDBClient("test_env", "test_tok")
	result := client.GetEnvironment()
	expected := "test_env"

	if result != expected {
		t.Errorf("GetEnvironment() = %s; want %s", result, expected)
	}
}
