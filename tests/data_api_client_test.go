package stragollum_test

import (
	"stragollum/pkg/stragollum"
	"testing"
)

func TestDataAPIClient(t *testing.T) {
	// Test case 1: Environment and Token are provided
	t.Run("WithEnvironmentAndToken", func(t *testing.T) {
		testEnv := stragollum.EnvironmentDev
		testTok := "test_tok_1"
		client := stragollum.NewDataAPIClient(&testEnv, &testTok)

		if client.Environment() != testEnv {
			t.Errorf("Environment() = %v; want %v", client.Environment(), testEnv)
		}
		if client.Token() == nil || *client.Token() != testTok {
			t.Errorf("Token() = %v; want %v", client.Token(), testTok)
		}
	})

	// Test case 2: Environment is provided, Token is nil
	t.Run("WithEnvironmentNoToken", func(t *testing.T) {
		testEnv := stragollum.EnvironmentTest
		client := stragollum.NewDataAPIClient(&testEnv, nil)

		if client.Environment() != testEnv {
			t.Errorf("Environment() = %v; want %v", client.Environment(), testEnv)
		}
		if client.Token() != nil {
			t.Errorf("Token() = %v; want nil", client.Token())
		}
	})

	// Test case 3: Environment is nil (should default to Prod), Token is provided
	t.Run("NoEnvironmentWithToken", func(t *testing.T) {
		testTok := "test_tok_2"
		client := stragollum.NewDataAPIClient(nil, &testTok)

		expectedEnv := stragollum.EnvironmentProd
		if client.Environment() != expectedEnv {
			t.Errorf("Environment() = %v; want %v", client.Environment(), expectedEnv)
		}
		if client.Token() == nil || *client.Token() != testTok {
			t.Errorf("Token() = %v; want %v", client.Token(), testTok)
		}
	})

	// Test case 4: Environment is nil (should default to Prod), Token is nil
	t.Run("NoEnvironmentNoToken", func(t *testing.T) {
		client := stragollum.NewDataAPIClient(nil, nil)

		expectedEnv := stragollum.EnvironmentProd
		if client.Environment() != expectedEnv {
			t.Errorf("Environment() = %v; want %v", client.Environment(), expectedEnv)
		}
		if client.Token() != nil {
			t.Errorf("Token() = %v; want nil", client.Token())
		}
	})
}
