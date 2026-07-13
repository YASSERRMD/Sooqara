package api

import (
	"context"
	"testing"
	"time"

	"github.com/yasserrmd/sooqara/internal/store"
)

func TestHandlerNilFields(t *testing.T) {
	h := NewHandler(nil, nil, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestHandlerStoresReferences(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	h := NewHandler(s, nil, nil)
	if h.store == nil {
		t.Error("expected store to be set")
	}
}

func TestContextDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = ctx.Err()
}
