package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

// App is cli wrapper that do some common operation and creates signal handler.
type App struct {
	Flags []cli.Flag
	Main  func(ctx context.Context)
}

func (a *App) wrapMain(c *cli.Context) error {
	ctx, cancel := context.WithCancel(c.Context)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("\nReceives signal: %v\n", sig)
		cancel()
	}()

	// Panic handling.
	defer func() {
		if r := recover(); r != nil {
			log.Println("Main recovered: ", r)
			debug.PrintStack()
		}
	}()

	a.Main(ctx)
	time.Sleep(3 * time.Second)
	log.Println("terminated")

	return nil
}

// Run setups everything and runs Main.
func (a *App) Run() {

	app := cli.NewApp()
	app.Flags = a.Flags
	app.Action = a.wrapMain

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
