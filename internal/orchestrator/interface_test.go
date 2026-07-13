package orchestrator

import (
	"testing"
)

func TestAllStagesImplemented(t *testing.T) {
	// Verify all 6 stages exist and implement the Stage interface
	var _ Stage = &analyseStage{}
	var _ Stage = &copyStage{}
	var _ Stage = &imageStage{}
	var _ Stage = &videoStage{}
	var _ Stage = &assembleStage{}
	var _ Stage = &noopStage{}
}
