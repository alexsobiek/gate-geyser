package main

import (
	"os"

	gategeyser "github.com/alexsobiek/gate-geyser"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	proxy.Plugins = append(
		proxy.Plugins,
		gategeyser.Plugin(".%s", getEnvOrDefault("GEYSER_LISTEN_ADDR", ":25566")),
	)

	gate.Execute()
}
