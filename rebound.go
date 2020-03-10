package rebound

import (
	"context"
	"time"

	"github.com/luukdegram/rebound/ecs"
	"github.com/luukdegram/rebound/internal/display"
	"github.com/luukdegram/rebound/internal/thread"
	"github.com/luukdegram/rebound/shaders"
)

//RunOptions allows you to set the initial width, height and title of the application
type RunOptions struct {
	Width  int
	Height int
	Title  string
}

// Run starts a new Rebound Application. It will initialize all base systems needed to run the engine.
// Those base systems include the display system and a basic rendering system.
func Run(options RunOptions, setup func()) error {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan error, 1)
	go func() {
		defer cancel()
		defer close(ch)
		if err := run(options, setup); err != nil {
			ch <- err
		}
	}()

	thread.Run(ctx)
	return <-ch
}

// run contains the main loop of the game.
// It handles updating of every system added to the manager
func run(options RunOptions, setup func()) error {
	window := display.Default()
	err := window.Init(options.Width, options.Height, options.Title)
	if err != nil {
		return err
	}

	setup()

	st := time.Now()
	for !window.ShouldClose() {
		delta := time.Now().Sub(st).Seconds() * 1000
		ecs.GetManager().Update(delta)
		window.Update()
		st = time.Now()
	}

	CleanUp()
	shaders.CleanUp()
	window.Close()

	return nil
}
