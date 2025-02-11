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
		gategeyser.Plugin(
			getEnvOrDefault("USERNAME_FORMAT", ".%s"),
			getEnvOrDefault("GEYSER_LISTEN_ADDR", "0.0.0.0:25566"),
			getEnvOrDefault("FLOODGATE_KEY_PATH", "/gate/floodgate.pem"),
		),
	)

	gate.Execute()
}
