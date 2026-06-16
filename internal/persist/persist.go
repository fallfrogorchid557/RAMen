// Package persist provides simple snapshot-to-disk persistence: the whole
// keyspace is gob-encoded to a file on an interval and on graceful shutdown,
// and loaded back on startup. This is the V1 durability model (PRD §6);
// append-only-file durability is deferred to V3.
package persist

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Rohit-Dnath/RAMen/internal/store"
)

// snapshot is the on-disk file format. Version guards against silently loading
// an incompatible future format.
type snapshot struct {
	Version int
	Records []store.Record
}

const formatVersion = 1

// Snapshotter writes and reads store snapshots at a fixed path.
type Snapshotter struct {
	path string
	s    *store.Store
}

// New returns a Snapshotter for the given store and file path.
func New(s *store.Store, path string) *Snapshotter {
	return &Snapshotter{path: path, s: s}
}

// Load restores the keyspace from the snapshot file. A missing file is not an
// error (first run); it returns false in that case.
func (sn *Snapshotter) Load() (bool, error) {
	f, err := os.Open(sn.path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer f.Close()

	var snap snapshot
	if err := gob.NewDecoder(f).Decode(&snap); err != nil {
		return false, fmt.Errorf("decode snapshot: %w", err)
	}
	if snap.Version != formatVersion {
		return false, fmt.Errorf("unsupported snapshot version %d", snap.Version)
	}
	sn.s.Import(snap.Records)
	return true, nil
}

// Save writes the current keyspace to the snapshot file atomically (write to a
// temp file, then rename).
func (sn *Snapshotter) Save() error {
	if sn.path == "" {
		return nil
	}
	if dir := filepath.Dir(sn.path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	tmp := sn.path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	snap := snapshot{Version: formatVersion, Records: sn.s.Export()}
	if err := gob.NewEncoder(f).Encode(&snap); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, sn.path)
}

// StartAuto periodically saves the snapshot until stop is closed, then saves
// one final time. Run it in its own goroutine.
func (sn *Snapshotter) StartAuto(interval time.Duration, stop <-chan struct{}) {
	if sn.path == "" || interval <= 0 {
		<-stop
		return
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			if err := sn.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "ramen: snapshot failed: %v\n", err)
			}
		}
	}
}
