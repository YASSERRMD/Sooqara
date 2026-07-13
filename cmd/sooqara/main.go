package main

import (
	"fmt"
	"os"

	"github.com/yasserrmd/sooqara/internal/config"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("config loaded: rpm=%d workers=%d addr=%s\n", cfg.RPM, cfg.Workers, cfg.Addr)
}
