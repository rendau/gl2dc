package rest

import (
	"context"
	"github.com/rendau/gl2dc/internal/domain/core"
	"net/http"
	"time"
)

type St struct {
	listen string
	cr     *core.St

	server *http.Server
	lChan  chan error
}

func New(listen string, cr *core.St) *St {
	api := &St{
		listen: listen,
		cr:     cr,
	}

	api.server = &http.Server{
		Addr:         listen,
		Handler:      api.createRouter(),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}

	return api
}

func (a *St) Start() {
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.lChan <- err
		}
	}()
}

func (a *St) Wait() <-chan error {
	return a.lChan
}

func (a *St) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
