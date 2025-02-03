package gategeyser

import (
	"fmt"
	"net"

	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/profile"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/uuid"
)

func (p *GateGeyserPlugin) getGeyserConnection(addr net.Addr) (*GeyserConnection, bool) {
	// Check if the connection is in the map
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, conn := range p.connections {
		if conn.RemoteAddr().String() == addr.String() {
			return conn, true
		}
	}

	return nil, false
}

func (p *GateGeyserPlugin) isGeyserConnection(addr net.Addr) bool {
	// Check if the connection is in the map
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, conn := range p.connections {
		if conn.RemoteAddr().String() == addr.String() {
			return true
		}
	}

	return false
}

func (p *GateGeyserPlugin) onPreLogin(e *proxy.PreLoginEvent) {
	// Check if the connection is in the map
	conn, ok := p.getGeyserConnection(e.Conn().RemoteAddr())

	if ok {
		hostname, data, err := p.floodgate.ReadHostname(e.Conn().VirtualHost().String())

		if err != nil || hostname == "" || data == nil {
			p.log.Info("Disconnecting player", "reason", "Failed to read hostname", "error", err)
			e.Deny(&component.Text{Content: "Failed to read bedrock hostname"})
			return
		}

		conn.BedrockData = data

		e.ForceOfflineMode()
	}
}

func (p *GateGeyserPlugin) onGameProfile(e *proxy.GameProfileRequestEvent) {
	conn, ok := p.getGeyserConnection(e.Conn().RemoteAddr())

	if ok {
		uid, err := conn.BedrockData.JavaUuid()

		if err != nil || uid == uuid.Nil {
			p.log.Info("Disconnecting player", "reason", "Failed to get UUID from XUID for game profile", "error", err)
			conn.Close()
			return
		}

		gameProfile := profile.GameProfile{
			Name: fmt.Sprintf(p.nameFormat, e.GameProfile().Name),
			ID:   uid,
		}

		skin, err := GetSkin(conn.Xuid)

		if err == nil && skin != nil {
			gameProfile.Properties = append(gameProfile.Properties, profile.Property{
				Name:      "textures",
				Value:     skin.Value,
				Signature: skin.Signature,
			})
		}

		linkedAccount, err := GetLinkedAccount(conn.Xuid)

		if linkedAccount != nil && linkedAccount.JavaID != uuid.Nil {
			// TODO: get skin for linked account
			p.log.Info("Bedrock player logged in as Java player", "bedrock", e.GameProfile().Name, "java", linkedAccount.JavaName)
			gameProfile.ID = linkedAccount.JavaID
			gameProfile.Name = linkedAccount.JavaName
		}

		e.SetGameProfile(gameProfile)
	}
}
