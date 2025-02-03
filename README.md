# Gate Geyser
This [Gate](https://gate.minekube.com/) plugin adds support for standalone [Geyser](https://geysermc.org/) proxies to be placed in front of your Gate instance

# Plugin Usage

```go
func main() {
  // Format for changing Bedrock usernames
  // to prevent conflicts with Java usernames
  usernameFormat := ".%s"
  // Listener address for Geyser to connect to
  // WARNING: You must configure a firewall to
  // only allow geyser to connect to this address
  // Non-geyser connections will bypass Mojang
  // authentication
  listenAddr := ":25566"

	proxy.Plugins = append(
		proxy.Plugins,
		gategeyser.Plugin(usernameFormat, listenAddr),
	)

	gate.Execute()
}
```

# Geyser Setup
1. Grab a copy of the [Geyser Standalone Jar](https://geysermc.org/download/)
2. Configure Geyser:
   - Set `remote.address` and `remote.port` to your Gate instance
   - Set `use-proxy-protocol` to true
3. Optionally setup [floodgate](https://geysermc.org/download/?project=floodgate) (highly recommended!)
   - Download floodgate on all of your backend servers
   - Use one key.pem for all (generate it on the first server, copy to all others)
   - Copy key.pem to your Geyser standalone
   - Set `remote.auth-type` to `floodgate` in your Geyser standalone's config

# Docker
A pre-made Gate proxy docker image w/geyser plugin is available. Use the `GEYSER_LISTEN_ADDR` environment variable to set
the address for the listener which Geyser Standalone will connect to. WARNING: You must properly configure a firewall
to prevent non-Geyser connections on this listener. Any other connections will result in a Mojang authentication bypass.

```sh
docker run \
	-e "GEYSER_LISTEN_ADDR=:25566" \
	-v /path/to/config.yaml:/gate/config.yaml \
   -v /path/to/floodgate/key.pem:/gate/floodgate.pem \
	-p 25565:25565 \
	ghcr.io/alexsobiek/gate-geyser:main
```
