package main

import (
	"fmt"

	gategeyser "github.com/alexsobiek/gate-geyser"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
	proxy.Plugins = append(
		proxy.Plugins,
		gategeyser.Plugin(".%s", []string{"172.30.1.0/24"}),
	)

	// gate.Execute()

	xuid, err := gategeyser.GetXuid("AlexSobiek")
	if err != nil {
		panic(err)
	}

	fmt.Println(xuid.Xuid)
}
