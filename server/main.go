package main

import (
	"log"
	"os"
	"strings"

	gategeyser "github.com/alexsobiek/gate-geyser"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type Settings struct {
	TrustedProxies []string
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getSettings() *Settings {
	proxyEnv := getEnvOrDefault("TRUSTED_PROXIES", "172.16.0.0/12")
	proxyList := strings.Split(proxyEnv, ",")

	return &Settings{
		TrustedProxies: proxyList,
	}
}

func main() {
	settings := getSettings()
	log.Println("Starting Gate with Geyser support")
	log.Printf("Trusted Proxies: %v\n", settings.TrustedProxies)

	proxy.Plugins = append(
		proxy.Plugins,
		gategeyser.Plugin(".%s", settings.TrustedProxies),
	)

	gate.Execute()
}
