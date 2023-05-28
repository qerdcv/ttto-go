package runner

import (
	"golang.org/x/sync/errgroup"

	"github.com/qerdcv/ttto/internal/ttto/runner/http"
)

type Runner interface {
	Run() error
}

var _ Runner = (*http.Runner)(nil)

func Run(runners ...Runner) error {
	errG := new(errgroup.Group)

	for _, r := range runners {
		errG.Go(r.Run)
	}

	return errG.Wait()
}
