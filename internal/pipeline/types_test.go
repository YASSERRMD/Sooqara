package pipeline

import (
	"testing"

	"github.com/yasserrmd/sooqara/internal/store"
)

func TestAllPipelineTypesDefined(t *testing.T) {
	// ProductAnalysis
	_ = ProductAnalysis{}
	_ = PhotoQuality{}
	// CopySet
	_ = CopySet{}
	// VideoConfig
	_ = VideoConfig{}
}

func TestArtifactKindConstants(t *testing.T) {
	if store.ArtifactAnalysis != "analysis" {
		t.Error("ArtifactAnalysis constant wrong")
	}
	if store.ArtifactCopy != "copy" {
		t.Error("ArtifactCopy constant wrong")
	}
	if store.ArtifactImage != "image" {
		t.Error("ArtifactImage constant wrong")
	}
	if store.ArtifactVideo != "video" {
		t.Error("ArtifactVideo constant wrong")
	}
}

func TestStateConstants(t *testing.T) {
	states := []store.State{
		store.StateQueued, store.StateAnalysing, store.StateCopywriting,
		store.StateImaging, store.StateVideoing, store.StateAssembling,
		store.StateDone, store.StateFailed, store.StateCancelled,
	}
	for _, s := range states {
		if s == "" {
			t.Error("state constant is empty")
		}
	}
}
