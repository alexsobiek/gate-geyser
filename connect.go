package gategeyser

import (
	"fmt"
	"net"

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
	_, ok := p.getGeyserConnection(e.Conn().RemoteAddr())

	if ok {
		e.ForceOfflineMode()
	}
}

func (p *GateGeyserPlugin) onGameProfile(e *proxy.GameProfileRequestEvent) {

	conn, ok := p.getGeyserConnection(e.Conn().RemoteAddr())

	if ok {
		xuid, err := GetXuid(e.GameProfile().Name)

		if err != nil {
			p.log.Info("Disconnecting player", "reason", "Failed to get XUID for game profile", "username", e.GameProfile().Name, "error", err)
			conn.Close()
			return
		}

		gameProfile := profile.GameProfile{}

		skin, err := GetSkin(xuid.Xuid)
		
		if err == nil && skin != nil {
			gameProfile.Properties = append(gameProfile.Properties, profile.Property{
				Name:      "textures",
				Value:     skin.Value,
				Signature: skin.Signature,
			})
		}

		linkedAccount, err := GetLinkedAccount(xuid.Xuid)

		if linkedAccount != nil && linkedAccount.JavaID != uuid.Nil {
			p.log.Info("Bedrock player logged in as Java player", "bedrock", e.GameProfile().Name, "java", linkedAccount.JavaName)
			gameProfile.ID = linkedAccount.JavaID
			gameProfile.Name = linkedAccount.JavaName
		} else {
			p.log.Info("Bedrock player logged in", "username", e.GameProfile().Name)

			uuid, err := xuid.Uuid()

			if err != nil {
				p.log.Info("Disconnecting player", "reason", "Failed to get UUID for game profile", "username", e.GameProfile().Name, "error", err)
				conn.Close()
				return
			}

			gameProfile.ID = uuid
			gameProfile.Name = fmt.Sprintf(p.nameFormat, e.GameProfile().Name)
		}

		e.SetGameProfile(gameProfile)
	}
}
