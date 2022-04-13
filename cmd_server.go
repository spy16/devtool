package main

import (
	"context"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/spy16/devtool/pkg/httpserver"
	"github.com/spy16/devtool/pkg/log"
)

func cmdServe(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "An HTTP server with multitude of services",
		Long: `
Modes:
* teapot: All request are handled by returning "418 I'm a teapot"
* request-echo: All requests are just echoed in the response
`,
	}

	var addr, mode string
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr", "a", ":8080", "Bind address for server")
	flags.StringVarP(&mode, "mode", "m", "echo", "Server mode to use")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		var h http.Handler
		switch mode {
		case "echo", "request_echo":
			h = httpserver.Echo()

		default:
			log.Warnf("unknown mode '%s', falling back to 'teapot'", mode)
			mode = "teapot"
			h = httpserver.TeaPot()
		}

		log.Infof("starting server in '%s' mode at %s...", mode, addr)
		if err := httpserver.Serve(ctx, addr, h); err != nil {
			log.Fatalf("http utility server exited with error: %v", err)
		}
		log.Infof("HTTP utility server exited gracefully")
	}

	return cmd
}
