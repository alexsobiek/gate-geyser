package gategeyser

import (
	"context"
	"net"
	"sync"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func Plugin(nameFormat string, addr string) proxy.Plugin {
	return proxy.Plugin{
		Name: "GateGeyserPlugin",
		Init: func(ctx context.Context, p *proxy.Proxy) error {
			pl := &GateGeyserPlugin{
				ctx:         ctx,
				log:         logr.FromContextOrDiscard(ctx),
				proxy:       p,
				nameFormat:  nameFormat,
				connections: make(map[net.Addr]*GeyserConnection),
			}

			return pl.init(addr)
		},
	}
}

type GateGeyserPlugin struct {
	ctx         context.Context
	log         logr.Logger
	proxy       *proxy.Proxy
	nameFormat  string
	connections map[net.Addr]*GeyserConnection
	mu          sync.RWMutex
}

func (p *GateGeyserPlugin) init(addr string) error {
	eventMgr := p.proxy.Event()

	event.Subscribe(eventMgr, 0, p.onPreLogin)
	event.Subscribe(eventMgr, 0, p.onGameProfile)

	go func() {
		err := p.listenAndServe(p.ctx, addr)

		if err != nil {
			panic(err)
		}
	}()
	return nil
}
