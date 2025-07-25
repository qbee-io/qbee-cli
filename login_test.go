package client

import (
	"context"
	"os"
	"testing"
)

func Test_LoginGetAuthenticatedClient_WithQBEE_TOKEN(t *testing.T) {
	// Save original environment values
	originalToken := os.Getenv("QBEE_TOKEN")
	originalBaseURL := os.Getenv("QBEE_BASEURL")
	originalEmail := os.Getenv("QBEE_EMAIL")
	originalPassword := os.Getenv("QBEE_PASSWORD")

	// Clean up environment after test
	defer func() {
		os.Setenv("QBEE_TOKEN", originalToken)
		os.Setenv("QBEE_BASEURL", originalBaseURL)
		os.Setenv("QBEE_EMAIL", originalEmail)
		os.Setenv("QBEE_PASSWORD", originalPassword)
	}()

	// Clear conflicting environment variables
	os.Unsetenv("QBEE_EMAIL")
	os.Unsetenv("QBEE_PASSWORD")

	t.Run("With QBEE_TOKEN set", func(t *testing.T) {
		testToken := "test-token-123"
		os.Setenv("QBEE_TOKEN", testToken)

		ctx := context.Background()
		client, err := LoginGetAuthenticatedClient(ctx)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if client == nil {
			t.Fatal("Expected client to be created")
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
		
		os.Setenv("QBEE_TOKEN", testToken)
		os.Setenv("QBEE_BASEURL", testBaseURL)

		ctx := context.Background()
		client, err := LoginGetAuthenticatedClient(ctx)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if client == nil {
			t.Fatal("Expected client to be created")
		}

		if client.GetAuthToken() != testToken {
			t.Fatalf("Expected auth token to be '%s', got '%s'", testToken, client.GetAuthToken())
		}

		if client.GetBaseURL() != testBaseURL {
			t.Fatalf("Expected base URL to be '%s', got '%s'", testBaseURL, client.GetBaseURL())
		}
	})

	t.Run("Without QBEE_TOKEN should fallback to config file", func(t *testing.T) {
		// Clear all authentication environment variables
		os.Unsetenv("QBEE_TOKEN")
		os.Unsetenv("QBEE_EMAIL")
		os.Unsetenv("QBEE_PASSWORD")
		os.Unsetenv("QBEE_BASEURL")

		ctx := context.Background()
		_, err := LoginGetAuthenticatedClient(ctx)

		// Should get an error trying to read config file since we don't have one
		if err == nil {
			t.Fatal("Expected error when no authentication method is available")
		}
	})
}