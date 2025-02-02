package gategeyser

import (
	"fmt"
	"net"

	"go.minekube.com/gate/pkg/edition/java/profile"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/uuid"
)

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
	if p.isGeyserConnection(e.Conn().RemoteAddr()) {
		e.ForceOfflineMode()
	}
}

func (p *GateGeyserPlugin) onGameProfile(e *proxy.GameProfileRequestEvent) {
	if p.isGeyserConnection(e.Conn().RemoteAddr()) {

		xuid, err := GetXuid(e.GameProfile().Name)

		if err != nil {
			p.log.Error(err, "Failed to get XUID for game profile")
			return
		}

		linkedAccount, err := GetLinkedAccount(xuid.Xuid)

		if linkedAccount != nil && linkedAccount.JavaID != uuid.Nil {
			p.log.Info("Bedrock player logged in as Java player", "bedrock", e.GameProfile().Name, "java", linkedAccount.JavaName)
			e.SetGameProfile(profile.GameProfile{
				ID:   linkedAccount.JavaID,
				Name: linkedAccount.JavaName,
			})
		} else {
			p.log.Info("Bedrock player logged in", "username", e.GameProfile().Name)

			uuid, err := xuid.Uuid()

			if err != nil {
				p.log.Error(err, "Failed to parse UUID from XUID")
				return
			}

			e.SetGameProfile(profile.GameProfile{
				ID:   uuid,
				Name: fmt.Sprintf(p.nameFormat, e.GameProfile().Name),
			})
		}
	}
}
