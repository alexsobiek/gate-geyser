package gategeyser

import (
	"fmt"
	"net"
	"strings"

	"github.com/pires/go-proxyproto"
	"go.minekube.com/gate/pkg/edition/java/profile"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/uuid"
)

type Conn struct {
	net.Conn
	cb func()
}

func (c *Conn) Close() error {
	c.cb()
	return c.Conn.Close()
}

func parseIP(ip string) (net.IP, error) {
	if strings.Contains(ip, ":") {
		var err error
		ip, _, err = net.SplitHostPort(ip)

		if err != nil {
			return nil, err
		}
	}

	return net.ParseIP(ip), nil
}

func (p *GateGeyserPlugin) isGeyserConnection(addr net.Addr) bool {
	// Check if the connection is in the map
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, ppConn := range p.connections {
		if ppConn.RemoteAddr().String() == addr.String() {
			return true
		}
	}

	return false
}

func (p *GateGeyserPlugin) onConnection(e *proxy.ConnectionEvent) {
	ip, err := parseIP(e.Connection().RemoteAddr().String())

	if err != nil {
		e.Connection().Close()
		panic(fmt.Errorf("Failed to handle connection event: %w", err))
	}

	if p.isTrustedProxy(ip) {
		// Use proxy protocol to get the real IP

		ppConn := proxyproto.NewConn(e.Connection())

		conn := &Conn{
			Conn: ppConn,
			cb: func() {
				p.mu.Lock()
				delete(p.connections, ppConn.RemoteAddr())
				p.mu.Unlock()
			},
		}

		e.SetConnection(conn)
		p.mu.Lock()
		p.connections[conn.RemoteAddr()] = ppConn
		p.mu.Unlock()
		return
	}
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
