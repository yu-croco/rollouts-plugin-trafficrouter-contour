package main

import (
	"os"

	"github.com/argoproj-labs/rollouts-contour-trafficrouter-plugin/pkg/plugin"

	rolloutsPlugin "github.com/argoproj/argo-rollouts/rollout/trafficrouting/plugin/rpc"
	goPlugin "github.com/hashicorp/go-plugin"
	"golang.org/x/exp/slog"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = goPlugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ARGO_ROLLOUTS_RPC_PLUGIN",
	MagicCookieValue: "trafficrouter",
}

func initLogger() {
	lvl := &slog.LevelVar{}
	lvl.Set(slog.LevelDebug)
	opts := slog.HandlerOptions{
		Level: lvl,
	}

	attrs := []slog.Attr{
		slog.String("plugin", "trafficrouter"),
		slog.String("vendor", "contour"),
	}
	opts.NewTextHandler(os.Stderr).WithAttrs(attrs)

	l := slog.New(opts.NewTextHandler(os.Stderr).WithAttrs(attrs))
	slog.SetDefault(l)
}

func main() {
	initLogger()
	rpcPluginImp := &plugin.RpcPlugin{}

	//  pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]goPlugin.Plugin{
		"RpcTrafficRouterPlugin": &rolloutsPlugin.RpcTrafficRouterPlugin{Impl: rpcPluginImp},
	}

	slog.Info("the plugin is running")
	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
