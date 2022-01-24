package app

import (
	"context"
	"sync"

	"github.com/hashicorp/go-multierror"
)

//App is the application interface, an object implementing this interface
//should follow the context protocol
type App interface {
	Run(ctx context.Context) error
}

//Runner is a parent app for apps to be run
type Runner struct {
	apps []App
}

//NewRunner returns *Runner
func NewRunner(apps ...App) *Runner {
	return &Runner{apps: apps}
}

// Add adds an App into collection of Apps to be run
func (r *Runner) Add(app App) *Runner {
	r.apps = append(r.apps, app)
	return r
}

//Run runs apps passed to the constructor concurrently,
//if one of them fails all the other running apps will be terminated
//by canceling context passed into their Run method
func (r *Runner) Run(ctx context.Context) error {
	var multiErrs error
	var mu sync.Mutex
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, a := range r.apps {
		wg.Add(1)
		app := a
		go func() {
			defer wg.Done()
			if err := app.Run(ctx); err != nil {
				mu.Lock()
				defer mu.Unlock()
				multiErrs = multierror.Append(multiErrs, err)
				cancel()
			}
		}()
	}

	wg.Wait()

	return multiErrs
}
