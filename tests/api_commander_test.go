package stragollum

import (
	"stragollum/pkg/stragollum"
	"testing"
)

func TestDataAPICommander(t *testing.T) {
	t.Run("WithToken", func(t *testing.T) {
		tok := "abc123"
		url := "https://api.example.com/api/json/v1/ks1"
		commander := stragollum.NewDataAPICommander(url, &tok)
		if commander.URL() != url {
			t.Errorf("URL() = %v; want %v", commander.URL(), url)
		}
		if commander.Token() == nil || *commander.Token() != tok {
			t.Errorf("Token() = %v; want %v", commander.Token(), tok)
		}
	})

	t.Run("NoToken", func(t *testing.T) {
		url := "https://api.example.com/api/json/v1/ks2"
		commander := stragollum.NewDataAPICommander(url, nil)
		if commander.URL() != url {
			t.Errorf("URL() = %v; want %v", commander.URL(), url)
		}
		if commander.Token() != nil {
			t.Errorf("Token() = %v; want nil", commander.Token())
		}
	})
}
