package api

import (
	"testing"
)

func Test_NewApiClient(t *testing.T) {
	client, err := NewApiClient()

	if err != nil {
		t.Fatalf("Err %+v creating ApiClient", err)
	}

	if client == nil {
		t.Fatal("Expected the client not to be nil")
	}
}
