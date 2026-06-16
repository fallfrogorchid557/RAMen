// Package dashboard serves a small read-only web UI: a live key browser plus
// hit/miss ratio and memory usage (PRD §7). It embeds a single static page and
// exposes a couple of JSON endpoints the page polls. It deliberately has no
// build step and no external assets.
package dashboard

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/Rohit-Dnath/RAMen/internal/server"
)

//go:embed static/index.html
var indexHTML []byte

// Dashboard wraps an HTTP server bound to a RAMen server's state.
type Dashboard struct {
	srv  *server.Server
	http *http.Server
}

// New builds a dashboard for srv listening on addr (e.g. ":8080").
func New(srv *server.Server, addr string) *Dashboard {
	d := &Dashboard{srv: srv}
	mux := http.NewServeMux()
	// Serve the single embedded page at the root. The page has no external
	// assets, so there is nothing else static to serve.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})
	mux.HandleFunc("/api/stats", d.handleStats)
	mux.HandleFunc("/api/keys", d.handleKeys)
	d.http = &http.Server{Addr: addr, Handler: mux}
	return d
}

// ListenAndServe runs the dashboard until ctx is cancelled.
func (d *Dashboard) ListenAndServe(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		d.http.Shutdown(shutCtx)
	}()
	if err := d.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (d *Dashboard) handleStats(w http.ResponseWriter, r *http.Request) {
	st := d.srv.Stats()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	writeJSON(w, map[string]any{
		"keys":              d.srv.Store().DBSize(),
		"connected_clients": st.Connections.Load(),
		"total_connections": st.TotalConns.Load(),
		"commands":          st.Commands.Load(),
		"cache_hits":        st.CacheHits.Load(),
		"cache_misses":      st.CacheMisses.Load(),
		"cache_hit_ratio":   st.HitRatio(),
		"memory_bytes":      mem.Alloc,
		"uptime_seconds":    int(d.srv.Uptime().Seconds()),
	})
}

func (d *Dashboard) handleKeys(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("match")
	if pattern == "" {
		pattern = "*"
	}
	keys := d.srv.Store().Keys(pattern)
	const limit = 500 // protect the browser from huge keyspaces
	truncated := false
	if len(keys) > limit {
		keys = keys[:limit]
		truncated = true
	}
	type keyInfo struct {
		Key  string `json:"key"`
		Type string `json:"type"`
	}
	infos := make([]keyInfo, 0, len(keys))
	for _, k := range keys {
		infos = append(infos, keyInfo{Key: k, Type: d.srv.Store().Type(k)})
	}
	writeJSON(w, map[string]any{"keys": infos, "truncated": truncated})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
