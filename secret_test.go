package client

import (
	"context"
	"regexp"
	"testing"
)

// secretRefPattern matches a secret reference, e.g. secret:ddb285d5-4f9c-4815-8fcf-7e6c4ba1d5ce (UUID4).
var secretRefPattern = regexp.MustCompile(
	`^secret:[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
)

func Test_CreateSecret(t *testing.T) {
	if !testingHasCredentials {
		t.Skip("Skipping test because QBEE_EMAIL and QBEE_PASSWORD are not set")
	}

	ctx := context.Background()

	cli, err := LoginGetAuthenticatedClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	secretRef, err := cli.CreateSecret(ctx, "qbee-cli-integration-test-secret")
	if err != nil {
		t.Fatal(err)
	}

	if !secretRefPattern.MatchString(secretRef) {
		t.Fatalf("expected secret_ref in format secret:<uuid4>, got %q", secretRef)
	}
}
