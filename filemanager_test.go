package client

import (
	"context"
	"os"
	"testing"
)

func Test_FileManager_Sync(t *testing.T) {
	// Create a new FileManager
	if os.Getenv("QBEE_EMAIL") == "" || os.Getenv("QBEE_PASSWORD") == "" {
		t.Skip("Skipping test because QBEE_EMAIL and QBEE_PASSWORD are not set")
	}

	ctx := context.Background()

	cli, err := LoginGetAuthenticatedClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	m := NewFileManager()
	// Set the client
	m.WithClient(cli)

	// Set the delete flag
	m.WithDelete(true)
	// Set the dryrun flag
	m.WithDryRun(false)

	if err := m.Sync(ctx, ".github", "/.github"); err != nil {
		t.Fatal(err)
	}
	cli.WithAuthToken("")

	if err := m.Sync(ctx, ".github", "/.github"); err != nil {
		t.Fatal(err)
	}

	// Synchronize the local directory with the FileManager

}
