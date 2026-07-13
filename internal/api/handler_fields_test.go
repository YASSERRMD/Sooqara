package api

import (
	"testing"

	"github.com/yasserrmd/sooqara/internal/store"
)

func TestHandlerIsEmpty(t *testing.T) {
	h := &Handler{}
	if h.store != nil || h.prov != nil || h.limiter != nil {
		t.Error("new handler should have nil fields")
	}
}

func TestHandlerNonEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	h := NewHandler(s, nil, nil)
	if h.store == nil {
		t.Error("store should not be nil")
	}
}
