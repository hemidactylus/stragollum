package stragollum_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"stragollum/pkg/stragollum"
	"testing"
)

func TestDatabase_DropCollection(t *testing.T) {
	const testCollectionName = "coll_drop_test"

	// Success case
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": {"ok": 1}}`)
	}))
	defer successServer.Close()

	token := "dummy"
	env := stragollum.EnvironmentProd
	client := stragollum.NewDataAPIClient(&env, &token)
	db := client.GetDatabase(successServer.URL, nil, "ks1")

	err := db.DropCollection(testCollectionName)
	if err != nil {
		t.Fatalf("DropCollection (success) failed: %v", err)
	}

	// Error: status.ok != 1
	nokServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status": {"ok": 0}}`)
	}))
	defer nokServer.Close()
	db = client.GetDatabase(nokServer.URL, nil, "ks1")
	err = db.DropCollection(testCollectionName)
	if err == nil {
		t.Error("Expected error when status.ok != 1, got nil")
	}

	// Error: missing status
	missingStatusServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"notstatus": {}}`)
	}))
	defer missingStatusServer.Close()
	db = client.GetDatabase(missingStatusServer.URL, nil, "ks1")
	err = db.DropCollection(testCollectionName)
	if err == nil {
		t.Error("Expected error when status is missing, got nil")
	}

	// Error: missing ok
	missingOkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status": {}}`)
	}))
	defer missingOkServer.Close()
	db = client.GetDatabase(missingOkServer.URL, nil, "ks1")
	err = db.DropCollection(testCollectionName)
	if err == nil {
		t.Error("Expected error when status.ok is missing, got nil")
	}
}
