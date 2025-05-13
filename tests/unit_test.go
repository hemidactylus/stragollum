package stragollum

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

func TestGetDatabase(t *testing.T) {
	clientToken := "client_token"
	client := stragollum.NewDataAPIClient(nil, &clientToken)
	apiEndpoint := "https://api.example.com"

	t.Run("AllParamsProvided", func(t *testing.T) {
		customToken := "custom_token"
		customKeyspace := "ks1"
		db := client.GetDatabase(apiEndpoint, &customToken, customKeyspace)
		if db == nil {
			t.Fatal("GetDatabase returned nil")
		}
		if db.Keyspace() != customKeyspace {
			t.Errorf("Keyspace() = %v; want %v", db.Keyspace(), customKeyspace)
		}
		if db.Token() == nil || *db.Token() != customToken {
			t.Errorf("Token = %v; want %v", db.Token(), customToken)
		}
		if db.ApiEndpoint() != apiEndpoint {
			t.Errorf("ApiEndpoint = %v; want %v", db.ApiEndpoint(), apiEndpoint)
		}
	})

	t.Run("TokenOmitted", func(t *testing.T) {
		customKeyspace := "ks2"
		db := client.GetDatabase(apiEndpoint, nil, customKeyspace)
		if db.Token() == nil || *db.Token() != clientToken {
			t.Errorf("Token = %v; want %v (client token)", db.Token(), clientToken)
		}
	})

	t.Run("KeyspaceOmitted", func(t *testing.T) {
		customToken := "custom_token2"
		db := client.GetDatabase(apiEndpoint, &customToken, "")
		if db.Keyspace() != "default_keyspace" {
			t.Errorf("Keyspace() = %v; want default_keyspace", db.Keyspace())
		}
	})

	t.Run("TokenAndKeyspaceOmitted", func(t *testing.T) {
		db := client.GetDatabase(apiEndpoint, nil, "")
		if db.Token() == nil || *db.Token() != clientToken {
			t.Errorf("Token = %v; want %v (client token)", db.Token(), clientToken)
		}
		if db.Keyspace() != "default_keyspace" {
			t.Errorf("Keyspace() = %v; want default_keyspace", db.Keyspace())
		}
	})
}
