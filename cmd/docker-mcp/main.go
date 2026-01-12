package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"

	"github.com/docker/mcp-gateway/cmd/docker-mcp/commands"
	"github.com/docker/mcp-gateway/cmd/docker-mcp/version"
	"github.com/docker/mcp-gateway/pkg/config"
	"github.com/docker/mcp-gateway/pkg/features"
	"github.com/docker/mcp-gateway/pkg/gateway"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Check gateway.yaml for skipDesktopCheck BEFORE features initialization
	// This must happen early because features.New() checks DOCKER_MCP_IN_CONTAINER
	if os.Getenv("DOCKER_MCP_IN_CONTAINER") != "1" {
		if data, err := config.ReadGatewayDefaults(); err == nil && len(data) > 0 {
			if defaults, err := gateway.ParseGatewayDefaults(data); err == nil && defaults.SkipDesktopCheck {
				os.Setenv("DOCKER_MCP_IN_CONTAINER", "1")
			}
		}
	}

	// We need to preserve CWD as paths.Init will change it.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if plugin.RunningStandalone() {
		os.Args = append([]string{os.Args[0], "mcp"}, os.Args[1:]...)
	}

	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		return commands.Root(ctx, cwd, dockerCli, features.New(ctx, dockerCli))
	},
		manager.Metadata{
			SchemaVersion:    "0.1.0",
			Vendor:           "Docker Inc.",
			Version:          version.Version,
			ShortDescription: "Docker MCP Plugin",
		},
	)
}
