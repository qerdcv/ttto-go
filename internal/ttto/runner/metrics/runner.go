package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/qerdcv/ttto/internal/conf"
)

type Runner struct {
	*http.Server
}

func New(cfg conf.Metrics) *Runner {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return &Runner{
		Server: &http.Server{
			Addr:    cfg.Addr,
			Handler: mux,
		},
	}
}

func (m *Runner) Run() error {
	if err := m.Server.ListenAndServe(); err != nil {
		return fmt.Errorf("metrics listen and serve: %w", err)
	}

	return nil
}
