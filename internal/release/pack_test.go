package release

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestPackTarballCreatesArchive(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, "base")
	if err := os.MkdirAll(base, 0755); err != nil {
		t.Fatal(err)
	}
	testFile := filepath.Join(base, "hello.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(tmp, "out.tar.gz")
	if err := PackTarball(archive, base, []string{"hello.txt"}); err != nil {
		t.Fatalf("PackTarball failed: %v", err)
	}

	info, err := os.Stat(archive)
	if err != nil {
		t.Fatalf("archive not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("expected non-empty archive")
	}
}

func TestPackTarballMissingFile(t *testing.T) {
	tmp := t.TempDir()
	err := PackTarball(filepath.Join(tmp, "out.tar.gz"), tmp, []string{"nonexistent.txt"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"hello.txt", "hello.txt"},
		{"path/to/file", "path_to_file"},
		{"../escape", ".._escape"},
	}
	for _, tt := range tests {
		got := SanitizeFilename(tt.input)
		if got != tt.expect {
			t.Errorf("SanitizeFilename(%q) = %q, want %q", tt.input, got, tt.expect)
		}
	}
}

func TestListFiles(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(tmp, "b.go"), []byte("b"), 0644)

	files, err := ListFiles(tmp, "*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 .txt file, got %d", len(files))
	}
}

func TestListFilesNoMatch(t *testing.T) {
	tmp := t.TempDir()
	files, err := ListFiles(tmp, "*.xyz")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestExtractTarballContent(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, "base")
	os.MkdirAll(base, 0755)
	os.WriteFile(filepath.Join(base, "readme.md"), []byte("# Hello"), 0644)

	archive := filepath.Join(tmp, "test.tar.gz")
	PackTarball(archive, base, []string{"readme.md"})

	f, _ := os.Open(archive)
	defer f.Close()
	gr, _ := gzip.NewReader(f)
	tr := tar.NewReader(gr)

	for h, err := tr.Next(); err != io.EOF; h, err = tr.Next() {
		if err != nil {
			t.Fatalf("tar read error: %v", err)
		}
		if h.Name == "readme.md" {
			break
		}
	}
}
