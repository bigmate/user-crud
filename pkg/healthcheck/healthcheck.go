package healthcheck

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"user-crud/pkg/app"

	"github.com/hashicorp/go-multierror"
)

type Resource interface {
	Ping(ctx context.Context) error
}

type healthCheck struct {
	port       string
	path       string
	timeout    time.Duration
	resources  []Resource
	httpServer *http.Server
}

func New(ctx context.Context, options ...Option) app.App {
	hc := defaultHealthCheck()

	for _, apply := range options {
		apply(hc)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(hc.path, hc.pingHandler())

	s := &http.Server{
		Addr:    ":" + hc.port,
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return ctx
		},
	}

	hc.httpServer = s

	return hc
}

func (hc *healthCheck) Run(ctx context.Context) error {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var multiErr error

	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		mu.Lock()
		defer mu.Unlock()
		if err := hc.httpServer.Shutdown(context.Background()); err != nil {
			multiErr = multierror.Append(multiErr, err)
		}
	}()

	if err := hc.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		mu.Lock()
		defer mu.Unlock()
		multiErr = multierror.Append(multiErr, err)
	}

	wg.Wait()

	return multiErr
}

func (hc *healthCheck) pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), hc.timeout)
		defer cancel()

		var mu sync.Mutex
		var wg sync.WaitGroup
		var multiErr error

		wg.Add(len(hc.resources))

		for _, res := range hc.resources {
			resource := res
			go func() {
				defer wg.Done()
				if err := resource.Ping(ctx); err != nil {
					mu.Lock()
					defer mu.Unlock()
					multiErr = multierror.Append(multiErr, err)
				}
			}()
		}

		wg.Wait()

		if ctx.Err() != nil {
			status(w, "ping timed out", false)
			return
		}

		if multiErr != nil {
			status(w, multiErr.Error(), false)
			return
		}

		status(w, "all resources are good", true)
	}
}

func status(w http.ResponseWriter, message string, ok bool) {
	st := struct {
		Message string `json:"message"`
		Ok      bool   `json:"ok"`
	}{
		Message: message,
		Ok:      ok,
	}
	bs, _ := json.Marshal(st)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	w.Write(bs)
}
