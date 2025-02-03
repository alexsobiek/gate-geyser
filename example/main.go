package main

import (
	gategeyser "github.com/alexsobiek/gate-geyser"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
	proxy.Plugins = append(
		proxy.Plugins,
		gategeyser.Plugin(".%s", ":25566", "/etc/gate/floodgate.pem"),
	)

	gate.Execute()
}
