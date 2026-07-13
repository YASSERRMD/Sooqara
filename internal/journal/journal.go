package journal

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

// Journal is an append-only activity journal backed by SQLite.
type Journal struct {
	db    *sql.DB
	ch    chan Activity
	dropMu sync.Mutex
	dropped uint64
	wg    sync.WaitGroup
	stopOnce sync.Once
}

// New creates a new Journal with a buffered channel of size bufSize.
func New(db *sql.DB, bufSize int) *Journal {
	j := &Journal{
		db: db,
		ch: make(chan Activity, bufSize),
	}
	j.wg.Add(1)
	go j.writer()
	return j
}

// Record adds an activity to the journal buffer.
// It does not block the caller. If the buffer is full, the entry is dropped.
func (j *Journal) Record(ctx context.Context, a Activity) error {
	select {
	case j.ch <- a:
		return nil
	default:
		j.dropMu.Lock()
		j.dropped++
		j.dropMu.Unlock()
		return fmt.Errorf("journal buffer full, entry dropped")
	}
}

// Flush drains the buffer for graceful shutdown.
func (j *Journal) Flush(ctx context.Context) error {
	j.dropMu.Lock()
	defer j.dropMu.Unlock()

	for {
		select {
		case a := <-j.ch:
			if j.db != nil {
				if err := j.insert(a); err != nil {
					return err
				}
			}
		default:
			return nil
		}
	}
}

// Close stops the writer goroutine and closes the database.
func (j *Journal) Close() error {
	j.stopOnce.Do(func() {
		close(j.ch)
	})
	j.wg.Wait()
	if j.db != nil {
		return j.db.Close()
	}
	return nil
}

// Dropped returns the number of dropped journal entries.
func (j *Journal) Dropped() uint64 {
	j.dropMu.Lock()
	defer j.dropMu.Unlock()
	return j.dropped
}
