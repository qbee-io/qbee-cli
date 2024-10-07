package client

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

var testingHasCredentials = os.Getenv("QBEE_EMAIL") != "" && os.Getenv("QBEE_PASSWORD") != ""

func Test_FileManager_Token_Refresh(t *testing.T) {
	if !testingHasCredentials {
		t.Skip("Skipping test because QBEE_EMAIL and QBEE_PASSWORD are not set")
	}

	ctx := context.Background()

	cli, err := LoginGetAuthenticatedClient(ctx)
	if err != nil {
		t.Fatal(err)
	}
	cli.authToken = "invalid"

	m := NewFileManager().WithClient(cli)
	if err := m.SnapshotRemote(ctx, "/"); err != nil {
		t.Fatal(err)
	}

}

func Test_FileManager_Exclude_Include(t *testing.T) {

	if !testingHasCredentials {
		t.Skip("Skipping test because QBEE_EMAIL and QBEE_PASSWORD are not set")
	}

	ctx := context.Background()

	cli, err := LoginGetAuthenticatedClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	m := NewFileManager().WithClient(cli).WithDryRun(true).WithExcludes("cmd/,.git").WithIncludes("cmd/filemanager.go")

	if err := m.SnapshotLocal("."); err != nil {
		t.Fatal(err)
	}

	list := m.GetLocalSnapshot()
	if len(list) == 0 {
		t.Fatal("should have 1 file")
	}

	for _, f := range list {
		if f.Path == "cmd/filemanager.go" {
			return
		}
	}
	t.Fatal("file not found")
}

func Test_FileManager_Upload_Download(t *testing.T) {
	if !testingHasCredentials {
		t.Skip("Skipping test because QBEE_EMAIL and QBEE_PASSWORD are not set")
	}

	ctx := context.Background()

	cli, err := LoginGetAuthenticatedClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	m := NewFileManager().WithClient(cli)

	tempDir := os.TempDir()
	tmpFile := filepath.Join(tempDir, "filemanager_test.txt")
	fileContents := []byte("filemanager_test")
	if err := os.WriteFile(tmpFile, fileContents, 0600); err != nil {
		t.Fatal(err)
	}

	if err := m.UploadFile(ctx, "/", tmpFile, true); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(tmpFile); err != nil {
		t.Fatal(err)
	}

	if err := m.DownloadFile(ctx, "/filemanager_test.txt", tmpFile); err != nil {
		t.Fatal(err)
	}

	downloadContents, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(downloadContents) != string(fileContents) {
		t.Fatalf("expected %s, got %s", fileContents, downloadContents)
	}
}

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

	if err := m.Sync(ctx, ".github", "/"); err != nil {
		t.Fatal(err)
	}

	if err := m.SnapshotRemote(ctx, "/"); err != nil {
		t.Fatal(err)
	}

	files := m.GetRemoteSnapshot()
	if len(files) == 0 {
		t.Fatal("no files found")
	}

	if err := m.Remove(ctx, "/", true); err != nil {
		t.Fatal(err)
	}

	if err := m.SnapshotRemote(ctx, "/"); err != nil {
		t.Fatal(err)
	}

	files = m.GetRemoteSnapshot()
	if len(files) != 0 {
		t.Fatal("files found")
	}

}
