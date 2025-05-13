package stragollum_test

import (
	"stragollum/pkg/stragollum"
	"testing"
)

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
