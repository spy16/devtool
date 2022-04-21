package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/spy16/devtool/pkg/config"
)

const appName = "devtool"

var (
	Version = "N/A"
	Commit  = "N/A"
	BuiltOn = "N/A"

	cfg appConfigs

	rootCmd = &cobra.Command{
		Use:     appName,
		Short:   "A toolbox for developers providing random useful utilities",
		Version: fmt.Sprintf("%s\ncommit: %s\nbuilt-on: %s", Version, Commit, BuiltOn),
	}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	rootCmd.PersistentFlags().StringP("config", "c", "", "Override config file")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		opts := []config.Option{
			config.WithEnv(),
			config.WithCobra(cmd),
		}

		if err := config.Load(&cfg, opts...); err != nil {
			fmt.Printf("failed to load configs: %v\n", err)
			os.Exit(1)
		}
	}

	rootCmd.AddCommand(
		cmdServe(ctx),
	)

	_ = rootCmd.Execute()
}

type appConfigs struct{}
