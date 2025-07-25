package client

import (
	"testing"
)

func Test_LoginGetAuthenticatedClient_With_QBEE_TOKEN(t *testing.T) {
	t.Run("With QBEE_TOKEN set", func(t *testing.T) {
		testToken := "test-token-123"
		t.Setenv("QBEE_BASEURL", "") // Ensure QBEE_BASEURL is not set
		t.Setenv("QBEE_TOKEN", testToken)
		client, err := LoginGetAuthenticatedClient(t.Context())
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if client.GetAuthToken() != testToken {
			t.Fatalf("Expected auth token to be '%s', got '%s'", testToken, client.GetAuthToken())
		}

		// Should use default base URL when QBEE_BASEURL is not set
		if client.GetBaseURL() != DefaultBaseURL {
			t.Fatalf("Expected base URL to be '%s', got '%s'", DefaultBaseURL, client.GetBaseURL())
		}
	})

	t.Run("With QBEE_TOKEN and QBEE_BASEURL set", func(t *testing.T) {
		testToken := "test-token-456"
		testBaseURL := "https://custom.qbee.io"
		
		t.Setenv("QBEE_TOKEN", testToken)
		t.Setenv("QBEE_BASEURL", testBaseURL)

		client, err := LoginGetAuthenticatedClient(t.Context())
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if client.GetAuthToken() != testToken {
			t.Fatalf("Expected auth token to be '%s', got '%s'", testToken, client.GetAuthToken())
		}

		if client.GetBaseURL() != testBaseURL {
			t.Fatalf("Expected base URL to be '%s', got '%s'", testBaseURL, client.GetBaseURL())
		}
	})
}