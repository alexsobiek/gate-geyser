# Gate Geyser
This [Gate](https://gate.minekube.com/) plugin adds support for standalone [Geyser](https://geysermc.org/) proxies to be placed in front of your Gate instance

# Plugin Usage

```go
func main() {
  # List of trusted proxies in CIDR notation
  # Be careful! Proxies within this range will bypass
  # Mojang's online authentication!
  trustedProxies := []string{"172.30.1.3/32"}

	proxy.Plugins = append(
		proxy.Plugins,
		gategeyser.Plugin(".%s", settings.TrustedProxies),
	)

	gate.Execute()
}
```

# Running
1. Grab a copy of the [Geyser Standalone Jar](https://geysermc.org/download/)
2. Configure Geyser:
   - Set `remote.address` and `remote.port` to your Gate instance
   - Set `use-proxy-protocol` to true
3. Optionally setup (floodgate)[https://geysermc.org/download/?project=floodgate] (highly recommended!)
   - Download floodgate on all of your backend servers
   - Use one key.pem for all (generate it on the first server, copy to all others)
   - Copy key.pem to your Geyser standalone
   - Set `remote.auth-type` to `floodgate` in your Geyser standalone's config
