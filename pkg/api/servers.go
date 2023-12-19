package api

import (
	"context"
	"sync"

	"glide/pkg/api/http"
)

type ServerManager struct {
	httpServer *http.Server
	shutdownWG *sync.WaitGroup
}

func NewServerManager(httpConfig *http.ServerConfig) (*ServerManager, error) {
	httpServer, err := http.NewServer(httpConfig)
	// TODO: init other servers like gRPC in future
	if err != nil {
		return nil, err
	}

	return &ServerManager{
		httpServer: httpServer,
		shutdownWG: &sync.WaitGroup{},
	}, nil
}

func (mgr *ServerManager) Start() {
	if mgr.httpServer != nil {
		mgr.shutdownWG.Add(1)

		go func() {
			defer mgr.shutdownWG.Done()

			// TODO: log the error
			err := mgr.httpServer.Run()

			println(err)
		}()
	}
}

func (mgr *ServerManager) Shutdown(ctx context.Context) error {
	var err error

	if mgr.httpServer != nil {
		err = mgr.httpServer.Shutdown(ctx)
	}

	mgr.shutdownWG.Wait()

	return err
}