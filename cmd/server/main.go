//nolint:gocognit,funlen,cyclop
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/StevenCyb/GoCLI/pkg/cli"
	"github.com/StevenCyb/ServMock/pkg/ini"
	"github.com/StevenCyb/ServMock/pkg/model"
	"github.com/StevenCyb/ServMock/pkg/server"
	"github.com/StevenCyb/ServMock/pkg/setup"
	"github.com/StevenCyb/ServMock/pkg/watcher"
)

const checkFileChangeInterval = 1000
const shutdownTimeout = 15 * time.Second

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	c := cli.New(
		cli.Name("ServMock"),
		cli.Banner(`
   _____                 __  __            _
  / ____|               |  \/  |          | |
 | (___   ___ _ ____   _| \  / | ___   ___| | __
  \___ \ / _ \ '__\ \ / / |\/| |/ _ \ / __| |/ /
  ____) |  __/ |   \ V /| |  | | (_) | (__|   <
 |_____/ \___|_|    \_/ |_|  |_|\___/ \___|_|\_\`),
		cli.Description("A REST service mocking tool."),
		cli.Version("0.1.0"),
		cli.Argument(
			"path",
			cli.Validate(regexp.MustCompile(`^.+\.ini$`)),
			cli.Description("Path to behavior config file."),
			cli.Option(
				"listen",
				cli.Description("Port to listen on for incoming requests."),
				cli.Short('l'),
				cli.Default(":3000"),
			),
			cli.Handler(
				func(ctx *cli.Context) error {
					path := ctx.GetArgument("path")
					if path == nil {
						return fmt.Errorf("invalid or missing path: %v", path)
					}
					if _, err := os.Stat(*path); err != nil {
						return fmt.Errorf("invalid or missing path: %s", *path)
					}
					listen := ctx.GetOption("listen")
					if listen == nil || !regexp.MustCompile(`^(localhost|127\.0\.0\.1)?:\d+$`).MatchString(*listen) {
						return fmt.Errorf("invalid or missing listen: %s", *listen)
					}

					logger.Info("Service mock listen", "listen", *listen, "path", *path)

					s := server.New(*listen, &model.BehaviorSet{})

					configErr := make(chan error, 1)
					watcherErr := make(chan error, 1)
					w := watcher.NewWatcher(*path, checkFileChangeInterval)
					w.RegisterListener(func(path string) {
						logger.Info("Configuration file changed", "path", path)
						file, err := os.Open(path)
						if err != nil {
							watcherErr <- fmt.Errorf("failed to open config file: %w", err)
							return
						}
						defer file.Close()
						sections, err := ini.Parse(file, true)
						if err != nil {
							watcherErr <- fmt.Errorf("failed to parse config file: %w", err)
							return
						}

						bs, err := setup.Build(sections)
						if err != nil {
							configErr <- fmt.Errorf("failed to build behavior set: %w", err)
							return
						}

						s.SetBehaviorSet(bs)
					})
					w.Start()
					defer w.Stop()

					serverError := s.Start()
					stop := make(chan os.Signal, 1)
					signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

					select {
					case err := <-serverError:
						log.Printf("Server error: %v", err)
					case err := <-watcherErr:
						log.Printf("Watcher error: %v", err)
					case err := <-configErr:
						log.Printf("Configuration error: %v", err)
					case sig := <-stop:
						log.Printf("Received shutdown signal: %v", sig)
					}

					log.Println("Server is shutting down...")

					shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
					defer cancel()

					if err := s.Shutdown(shutdownCtx); err != nil {
						return err
					}

					log.Println("Server exited properly")

					return nil
				},
			),
		),
	)

	_, err := c.RunWith(os.Args)
	if err != nil {
		logger.Error("Error running CLI", "error", err)
		c.PrintHelp()
		os.Exit(1)
	}
}
