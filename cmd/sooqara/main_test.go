package main

import (
	"os/exec"
	"testing"
)

func TestMainExitsWithConfigError(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/sooqara")
	cmd.Env = append(os.Environ(), "AGNES_API_KEY=")
	out, _ := cmd.CombinedOutput()
	if len(out) == 0 {
		t.Fatal("expected config error output")
	}
}
