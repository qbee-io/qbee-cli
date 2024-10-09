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

	cli.WithAuthToken("invalid")

	m := NewFileManager().WithClient(cli)
	if err := m.SnapshotRemote(ctx, "/"); err != nil {
		t.Fatal(err)
	}

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

	testDir := createTestDirectoryStructure(t)

	if err := m.Sync(ctx, testDir, "/testDir"); err != nil {
		t.Fatal(err)
	}

	if err := m.Sync(ctx, testDir, "/"); err != nil {
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

func Test_FileManager_Exclude(t *testing.T) {

	if !testingHasCredentials {
		t.Skip("Skipping test because QBEE_EMAIL and QBEE_PASSWORD are not set")
	}

	ctx := context.Background()

	cli, err := LoginGetAuthenticatedClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	testDir := createTestDirectoryStructure(t)

	m := NewFileManager().WithClient(cli).WithDryRun(true).WithExcludes("subdir")

	if err := m.SnapshotLocal(testDir); err != nil {
		t.Fatal(err)
	}

	list := m.GetLocalSnapshot()

	tt := []struct {
		path        string
		shouldExist bool
	}{
		{"testdir/subdir", false},
		{"testdir", true},
		{"testdir/subdir/testfile2.txt", false},
		{"testdir/subdir/testfile3.txt", false},
	}

	for _, tc := range tt {
		if _, ok := list[tc.path]; ok != tc.shouldExist {
			t.Fatalf("expected %s to exist: %v", tc.path, tc.shouldExist)
		}
	}

}

func Test_FileManager_Exlude_Include(t *testing.T) {
	if !testingHasCredentials {
		t.Skip("Skipping test because QBEE_EMAIL and QBEE_PASSWORD are not set")
	}

	ctx := context.Background()

	cli, err := LoginGetAuthenticatedClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	testDir := createTestDirectoryStructure(t)

	m := NewFileManager().
		WithClient(cli).
		WithDryRun(true).
		WithExcludes("subdir").
		WithIncludes("subdir/testfile2.txt")

	if err := m.SnapshotLocal(testDir); err != nil {
		t.Fatal(err)
	}

	list := m.GetLocalSnapshot()
	if len(list) == 0 {
		t.Fatal("should have 1 file")
	}

	tt := []struct {
		path        string
		shouldExist bool
	}{
		{"testdir/subdir", true},
		{"testdir", true},
		{"testdir/subdir/testfile2.txt", true},
		{"testdir/subdir/testfile3.txt", false},
	}

	for _, tc := range tt {
		if _, ok := list[tc.path]; ok != tc.shouldExist {
			t.Fatalf("expected %s to exist: %v", tc.path, tc.shouldExist)
		}
	}
}

func createTestDirectoryStructure(t *testing.T) string {

	tmpDir := t.TempDir()

	testDir := filepath.Join(tmpDir, "testdir")

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(testDir, "testfile.txt"), []byte("testfile"), 0600); err != nil {
		t.Fatal(err)
	}

	testSubDir := filepath.Join(testDir, "subdir")

	if err := os.MkdirAll(testSubDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(testSubDir, "testfile2.txt"), []byte("testfile2"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(testSubDir, "testfile3.txt"), []byte("testfile2"), 0600); err != nil {
		t.Fatal(err)
	}

	return tmpDir
}
