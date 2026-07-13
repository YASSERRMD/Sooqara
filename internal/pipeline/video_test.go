package pipeline

import (
	"context"
	"os"
	"testing"

	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/store"
)

func TestCreateVideoJobPersistsVideoID(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	db, _ := store.Open(dbPath)
	defer db.Close()
	blob := store.NewFilesystemBlob(tmpDir + "/storage")
	s := store.NewStore(db, blob)

	job := store.NewJob(tmpDir+"/src.jpg", "test", "", 2)
	store.CreateJob(s.DB, job)

	// Create an image artifact first
	imgPath := tmpDir + "/img.jpg"
	os.WriteFile(imgPath, []byte("img"), 0644)
	a := store.NewArtifact(job.ID, store.ArtifactImage, 0)
	a.Path = &imgPath
	store.CreateArtifact(s.DB, a)

	// Set job to imaging state first
	store.Transition(s.DB, job.ID, store.StateQueued, store.StateAnalysing)
	store.Transition(s.DB, job.ID, store.StateAnalysing, store.StateCopywriting)
	store.Transition(s.DB, job.ID, store.StateCopywriting, store.StateImaging)

	analysis := &ProductAnalysis{ProductName: "X", Category: "C", ShapeDescription: "D"}

	fake := &FakeProvider{
		CreateVid: func(ctx context.Context, req provider.VideoRequest) (provider.VideoJob, error) {
			return provider.VideoJob{VideoID: "vid-123"}, nil
		},
	}

	_, err := CreateVideoJob(context.Background(), fake, s, job, analysis)
	if err != nil {
		t.Fatalf("CreateVideoJob failed: %v", err)
	}
}
