package stragollum_test

import (
	"fmt" // Added for URL construction
	"stragollum/pkg/stragollum"
	"testing"
)

func TestGetDatabase(t *testing.T) {
	clientToken := "client_token"
	var testEnv stragollum.Environment = stragollum.EnvironmentProd // Or some default
	client := stragollum.NewDataAPIClient(&testEnv, &clientToken)
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
			t.Errorf("Token() = %v; want %v", db.Token(), &customToken)
		}
		if db.ApiEndpoint() != apiEndpoint {
			t.Errorf("ApiEndpoint() = %v; want %v", db.ApiEndpoint(), apiEndpoint)
		}

		// Commander tests
		commander := db.Commander()
		if commander == nil {
			t.Fatal("Commander() returned nil")
		}
		expectedCommanderURL := fmt.Sprintf("%s/api/json/v1/%s", apiEndpoint, customKeyspace)
		if commander.URL() != expectedCommanderURL {
			t.Errorf("Commander().URL() = %v; want %v", commander.URL(), expectedCommanderURL)
		}
		if commander.Token() == nil || *commander.Token() != customToken {
			t.Errorf("Commander().Token() = %v; want %v", commander.Token(), &customToken)
		}
	})

	t.Run("TokenOmitted", func(t *testing.T) {
		customKeyspace := "ks2"
		db := client.GetDatabase(apiEndpoint, nil, customKeyspace)
		if db == nil {
			t.Fatal("GetDatabase returned nil")
		}
		if db.Token() == nil || *db.Token() != clientToken {
			t.Errorf("Token() = %v; want %v (client token)", db.Token(), &clientToken)
		}
		if db.Keyspace() != customKeyspace {
			t.Errorf("Keyspace() = %v; want %v", db.Keyspace(), customKeyspace)
		}

		// Commander tests
		commander := db.Commander()
		if commander == nil {
			t.Fatal("Commander() returned nil")
		}
		expectedCommanderURL := fmt.Sprintf("%s/api/json/v1/%s", apiEndpoint, customKeyspace)
		if commander.URL() != expectedCommanderURL {
			t.Errorf("Commander().URL() = %v; want %v", commander.URL(), expectedCommanderURL)
		}
		if commander.Token() == nil || *commander.Token() != clientToken {
			t.Errorf("Commander().Token() = %v; want %v (client token)", commander.Token(), &clientToken)
		}
	})

	t.Run("KeyspaceOmitted", func(t *testing.T) {
		customToken := "custom_token2"
		db := client.GetDatabase(apiEndpoint, &customToken, "")
		if db == nil {
			t.Fatal("GetDatabase returned nil")
		}
		if db.Keyspace() != stragollum.DefaultKeyspace { // Use the actual constant
			t.Errorf("Keyspace() = %v; want %v", db.Keyspace(), stragollum.DefaultKeyspace)
		}
		if db.Token() == nil || *db.Token() != customToken {
			t.Errorf("Token() = %v; want %v", db.Token(), &customToken)
		}

		// Commander tests
		commander := db.Commander()
		if commander == nil {
			t.Fatal("Commander() returned nil")
		}
		expectedCommanderURL := fmt.Sprintf("%s/api/json/v1/%s", apiEndpoint, stragollum.DefaultKeyspace)
		if commander.URL() != expectedCommanderURL {
			t.Errorf("Commander().URL() = %v; want %v", commander.URL(), expectedCommanderURL)
		}
		if commander.Token() == nil || *commander.Token() != customToken {
			t.Errorf("Commander().Token() = %v; want %v", commander.Token(), &customToken)
		}
	})

	t.Run("TokenAndKeyspaceOmitted", func(t *testing.T) {
		db := client.GetDatabase(apiEndpoint, nil, "")
		if db == nil {
			t.Fatal("GetDatabase returned nil")
		}
		if db.Token() == nil || *db.Token() != clientToken {
			t.Errorf("Token() = %v; want %v (client token)", db.Token(), &clientToken)
		}
		if db.Keyspace() != stragollum.DefaultKeyspace { // Use the actual constant
			t.Errorf("Keyspace() = %v; want %v", db.Keyspace(), stragollum.DefaultKeyspace)
		}

		// Commander tests
		commander := db.Commander()
		if commander == nil {
			t.Fatal("Commander() returned nil")
		}
		expectedCommanderURL := fmt.Sprintf("%s/api/json/v1/%s", apiEndpoint, stragollum.DefaultKeyspace)
		if commander.URL() != expectedCommanderURL {
			t.Errorf("Commander().URL() = %v; want %v", commander.URL(), expectedCommanderURL)
		}
		if commander.Token() == nil || *commander.Token() != clientToken {
			t.Errorf("Commander().Token() = %v; want %v (client token)", commander.Token(), &clientToken)
		}
	})
}
