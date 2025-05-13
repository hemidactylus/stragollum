//go:build integration
// +build integration

package stragollum_test

import (
	"testing"
)

func TestIntegration(t *testing.T) {
	// TODO: Add integration test logic here
	t.Run("DummyFailer", func(t *testing.T) {
		t.Errorf("Boo!")
	})
}
