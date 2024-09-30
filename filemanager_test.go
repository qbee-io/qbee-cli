package client

import (
	"context"
	"os"
	"testing"
)

var testingHasCredentials = os.Getenv("QBEE_EMAIL") != "" && os.Getenv("QBEE_PASSWORD") != ""

func Test_FileManager_Sync(t *testing.T) {
	// Create a new FileManager
	if !testingHasCredentials {
		t.Skip("Skipping test because QBEE_EMAIL and QBEE_PASSWORD are not set")
	}

	ctx := context.Background()

	cli, err := LoginGetAuthenticatedClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	m := NewFileManager().WithClient(cli).WithDelete(true)

	if err := m.Sync(ctx, ".github", "/.github"); err != nil {
		t.Fatal(err)
	}

	cli.WithAuthToken("invalid-test-refresh-token")

	if err := m.Sync(ctx, ".github", "/"); err != nil {
		t.Fatal(err)
	}

	if err := m.SnapshotRemote(ctx, "/"); err != nil {
		t.Fatal(err)
	}
	cli.WithAuthToken("invalid-test-refresh-token")

	files := m.GetRemoteSnapshot()
	if len(files) == 0 {
		t.Fatal("no files found")
	}

	if err := m.Remove(ctx, "/", true); err != nil {
		t.Fatal(err)
	}

	if err := m.SnapshotRemote(ctx, "/"); err == nil {
		t.Fatal("expected error")
	}

	files = m.GetRemoteSnapshot()
	if len(files) != 0 {
		t.Fatal("files found")
	}

	// Synchronize the local directory with the FileManager
}
