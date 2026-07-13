package pipeline

import (
	"testing"
)

func TestProductAnalysisZeroValue(t *testing.T) {
	var a ProductAnalysis
	if a.ProductName != "" {
		t.Error("zero-value ProductName should be empty")
	}
}

func TestPhotoQualityZeroValue(t *testing.T) {
	var pq PhotoQuality
	if pq.Lightning != "" {
		t.Error("zero-value PhotoQuality fields should be empty")
	}
}

func TestCopySetZeroValue(t *testing.T) {
	var cs CopySet
	if cs.Title != "" {
		t.Error("zero-value CopySet fields should be empty")
	}
}
