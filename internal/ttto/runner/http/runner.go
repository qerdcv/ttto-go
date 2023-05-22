package http

import (
	"fmt"

	"github.com/qerdcv/ttto/internal/conf"
	"github.com/qerdcv/ttto/internal/ttto/service"
)

type Runner struct {
	server *server
	cfg    conf.HTTP
}

func New(svc *service.Service, cfg conf.HTTP) *Runner {
	return &Runner{
		server: newServer(svc),
		cfg:    cfg,
	}
}

func (r *Runner) Run() error {
	fmt.Println(r.cfg.Addr)
	return r.server.Run(r.cfg.Addr)
}
