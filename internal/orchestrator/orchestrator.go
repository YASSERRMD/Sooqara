package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yasserrmd/sooqara/internal/journal"
	"github.com/yasserrmd/sooqara/internal/provider"
	"github.com/yasserrmd/sooqara/internal/provider/agnes"
	"github.com/yasserrmd/sooqara/internal/pipeline"
	"github.com/yasserrmd/sooqara/internal/store"
)

// Config holds orchestrator settings.
type Config struct {
	WorkerCount int
	PollInterval time.Duration
	GracePeriod  time.Duration
}

// DefaultConfig returns standard settings.
func DefaultConfig() Config {
	return Config{
		WorkerCount:  3,
		PollInterval: 2 * time.Second,
		GracePeriod:  60 * time.Second,
	}
}

// Orchestrator drives jobs through the pipeline.
type Orchestrator struct {
	cfg    Config
	store  *store.Store
	prov   provider.Provider
	journal *journal.Journal
	stop   chan struct{}
	wg     sync.WaitGroup
}

// New creates a new orchestrator.
func New(cfg Config, s *store.Store, prov provider.Provider, j *journal.Journal) *Orchestrator {
	if cfg.WorkerCount <= 0 {
		cfg.WorkerCount = 3
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 2 * time.Second
	}
	return &Orchestrator{
		cfg:     cfg,
		store:   s,
		prov:    prov,
		journal: j,
		stop:    make(chan struct{}),
	}
}

// Run starts the worker pool. Blocks until Stop is called.
func (o *Orchestrator) Run(ctx context.Context) error {
	for w := 0; w < o.cfg.WorkerCount; w++ {
		o.wg.Add(1)
		go o.worker(ctx, w)
	}

	// Start video poller in background
	videoCfg := pipeline.DefaultVideoConfig()
	poller := pipeline.NewVideoPoller(o.prov, o.store, videoCfg)
	go poller.Run(ctx)

	<-ctx.Done()
	o.stopWorkers(ctx)
	poller.Stop()
	return o.journal.Flush(ctx)
}

func (o *Orchestrator) worker(ctx context.Context, id int) {
	defer o.wg.Done()
	for {
		select {
		case <-o.stop:
			return
		default:
		}

		job, err := store.ClaimNext(o.store.DB, []store.State{
			store.StateQueued, store.StateAnalysing,
			store.StateCopywriting, store.StateImaging,
			store.StateAssembling,
		})
		if err != nil {
			time.Sleep(o.cfg.PollInterval)
			continue
		}
		if job == nil {
			time.Sleep(o.cfg.PollInterval)
			continue
		}

		stage := stageFor(job.State)
		stageCtx, cancel := context.WithTimeout(ctx, stage.Timeout())
		err = stage.Run(stageCtx, o.prov, o.store, job)
		cancel()

		if err != nil {
			o.handleStageError(job, err)
			continue
		}

		next := store.NextState(job.State)
		if err := store.Transition(o.store.DB, job.ID, job.State, next); err != nil {
			o.handleStageError(job, fmt.Errorf("transition %s -> %s: %w", job.State, next, err))
		}
	}
}

func (o *Orchestrator) stopWorkers(ctx context.Context) {
	close(o.stop)
	done := make(chan struct{})
	go func() { o.wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(o.cfg.GracePeriod):
	}
}

func (o *Orchestrator) handleStageError(job *store.Job, err error) {
	errStr := err.Error()
	o.store.DB.Exec("UPDATE jobs SET error = ?, state = 'failed' WHERE id = ?", errStr, job.ID)
}

func stageFor(state store.State) Stage {
	switch state {
	case store.StateQueued:
		return &analyseStage{}
	case store.StateAnalysing:
		return &copyStage{}
	case store.StateCopywriting:
		return &imageStage{}
	case store.StateImaging:
		return &videoStage{}
	case store.StateAssembling:
		return &assembleStage{}
	default:
		return &noopStage{}
	}
}

// Stage is a pipeline stage.
type Stage interface {
	Name() string
	Run(ctx context.Context, p provider.Provider, s *store.Store, job *store.Job) error
	Timeout() time.Duration
}

type analyseStage struct{}

func (s *analyseStage) Name() string          { return "analyse" }
func (s *analyseStage) Timeout() time.Duration { return 120 * time.Second }
func (s *analyseStage) Run(ctx context.Context, p provider.Provider, sto *store.Store, job *store.Job) error {
	arts, err := store.GetArtifactsByJob(sto.DB, job.ID)
	if err != nil {
		return err
	}
	for _, a := range arts {
		if a.Kind == store.ArtifactAnalysis && a.Payload != nil {
			return nil // already done
		}
	}
	return fmt.Errorf("analyse stage: not implemented")
}

type copyStage struct{}

func (s *copyStage) Name() string          { return "copy" }
func (s *copyStage) Timeout() time.Duration { return 120 * time.Second }
func (s *copyStage) Run(ctx context.Context, p provider.Provider, sto *store.Store, job *store.Job) error {
	return fmt.Errorf("copy stage: not implemented")
}

type imageStage struct{}

func (s *imageStage) Name() string          { return "image" }
func (s *imageStage) Timeout() time.Duration { return 180 * time.Second }
func (s *imageStage) Run(ctx context.Context, p provider.Provider, sto *store.Store, job *store.Job) error {
	return fmt.Errorf("image stage: not implemented")
}

type videoStage struct{}

func (s *videoStage) Name() string          { return "video" }
func (s *videoStage) Timeout() time.Duration { return 60 * time.Second }
func (s *videoStage) Run(ctx context.Context, p provider.Provider, sto *store.Store, job *store.Job) error {
	return fmt.Errorf("video stage: not implemented")
}

type assembleStage struct{}

func (s *assembleStage) Name() string          { return "assemble" }
func (s *assembleStage) Timeout() time.Duration { return 30 * time.Second }
func (s *assembleStage) Run(ctx context.Context, p provider.Provider, sto *store.Store, job *store.Job) error {
	return nil
}

type noopStage struct{}

func (s *noopStage) Name() string          { return "noop" }
func (s *noopStage) Timeout() time.Duration { return 10 * time.Second }
func (s *noopStage) Run(ctx context.Context, p provider.Provider, sto *store.Store, job *store.Job) error {
	return nil
}
