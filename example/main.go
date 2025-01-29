package main

import (
	gategeyser "github.com/alexsobiek/gate-geyser"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
	proxy.Plugins = append(
		proxy.Plugins,
		gategeyser.Plugin(".%s", []string{"172.30.1.0/24"}),
	)

	gate.Execute()
}
