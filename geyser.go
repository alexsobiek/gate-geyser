package gategeyser

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"github.com/yl2chen/cidranger"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Plugin(nameFormat string, trustedProxies []string) proxy.Plugin {
	return proxy.Plugin{
		Name: "GateGeyserPlugin",
		Init: func(ctx context.Context, p *proxy.Proxy) error {
			pl := &GateGeyserPlugin{
				ctx:            ctx,
				log:            logr.FromContextOrDiscard(ctx),
				proxy:          p,
				trustedProxies: cidranger.NewPCTrieRanger(),
				nameFormat:     nameFormat,
				connections:    make(map[net.Addr]*Conn),
			}

			for _, trustedProxy := range trustedProxies {
				_, network, err := net.ParseCIDR(trustedProxy)
				if err != nil {
					return fmt.Errorf("invalid trusted proxy: %w", err)
				}
				pl.trustedProxies.Insert(cidranger.NewBasicRangerEntry(*network))
			}

			return pl.init()
		},
	}
}

type GateGeyserPlugin struct {
	ctx            context.Context
	log            logr.Logger
	proxy          *proxy.Proxy
	trustedProxies cidranger.Ranger
	nameFormat     string
	connections    map[net.Addr]*Conn
	mu             sync.RWMutex
}

func (p *GateGeyserPlugin) init() error {

	eventMgr := p.proxy.Event()

	event.Subscribe(eventMgr, 0, p.onConnection)
	event.Subscribe(eventMgr, 0, p.onPreLogin)
	event.Subscribe(eventMgr, 0, p.onGameProfile)

	return nil
}

func (p *GateGeyserPlugin) isTrustedProxy(ip net.IP) bool {
	var ok bool
	var err error
	if ok, err = p.trustedProxies.Contains(ip); err != nil {
		return false
	}

	return ok
}
