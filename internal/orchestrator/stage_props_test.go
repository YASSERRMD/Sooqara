package orchestrator

import (
	"testing"
)

func TestAllStagesHavePositiveTimeout(t *testing.T) {
	stages := []Stage{&analyseStage{}, &copyStage{}, &imageStage{}, &videoStage{}, &assembleStage{}, &noopStage{}}
	for _, s := range stages {
		if s.Timeout() <= 0 {
			t.Errorf("%s: timeout = %v, want > 0", s.Name(), s.Timeout())
		}
	}
}

func TestAllStagesHaveUniqueNames(t *testing.T) {
	names := make(map[string]int)
	stages := []Stage{&analyseStage{}, &copyStage{}, &imageStage{}, &videoStage{}, &assembleStage{}, &noopStage{}}
	for _, s := range stages {
		names[s.Name()]++
	}
	for name, count := range names {
		if count > 1 {
			t.Errorf("duplicate stage name: %s", name)
		}
	}
}
