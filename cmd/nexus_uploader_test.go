package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/schollz/progressbar/v3"
)

func TestUploadSequentially(t *testing.T) {
	tempDir := t.TempDir()
	createTestFiles(t, tempDir, []string{"file1.txt", "file2.txt", "file3.txt"})

	files, err := collectFiles(tempDir)
	if err != nil {
		t.Fatalf("failed to collect files: %v", err)
	}

	bar := progressbar.Default(int64(len(files))) // Create a progress bar for testing

	err = uploadSequentially(files, bar)
	if err != nil {
		t.Errorf("uploadSequentially failed: %v", err)
	}
}

func TestUploadConcurrently(t *testing.T) {
	tempDir := t.TempDir()
	createTestFiles(t, tempDir, []string{"file1.txt", "file2.txt", "file3.txt"})

	files, err := collectFiles(tempDir)
	if err != nil {
		t.Fatalf("failed to collect files: %v", err)
	}

	bar := progressbar.Default(int64(len(files))) // Create a progress bar for testing

	concurrency = 2
	err = uploadConcurrently(files, bar)
	if err != nil {
		t.Errorf("uploadConcurrently failed: %v", err)
	}
}

func TestCollectFiles(t *testing.T) {
	tempDir := t.TempDir()
	createTestFiles(t, tempDir, []string{"file1.txt", "file2.txt", "subdir/file3.txt"})

	files, err := collectFiles(tempDir)
	if err != nil {
		t.Fatalf("collectFiles failed: %v", err)
	}

	expectedCount := 3
	if len(files) != expectedCount {
		t.Errorf("file count mismatch: got %d, want %d", len(files), expectedCount)
	}
}

func createTestFiles(t *testing.T, root string, paths []string) {
	t.Helper()

	for _, p := range paths {
		path := filepath.Join(root, p)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}

		file, err := os.Create(path)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		file.Close()
	}
}

func collectFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
