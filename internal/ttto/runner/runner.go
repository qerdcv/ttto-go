package runner

import "github.com/qerdcv/ttto/internal/ttto/runner/http"

type Runner interface {
	Run() error
}

var _ Runner = (*http.Runner)(nil)
