package gategeyser

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/alexsobiek/gate-geyser/floodgate"
	"github.com/pires/go-proxyproto"
	"go.minekube.com/gate/pkg/util/errs"
)

type GeyserConnection struct {
	net.Conn
	*floodgate.BedrockData
	closeCb func()
}

func (c *GeyserConnection) Close() error {
	c.closeCb()
	return c.Conn.Close()
}

func (p *GateGeyserPlugin) listenAndServe(ctx context.Context, addr string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer func() { _ = ln.Close() }()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() { <-ctx.Done(); _ = ln.Close() }()

	defer p.log.Info("stopped listening for new bedrock connections", "addr", addr)
	p.log.Info("listening for bedrock connections", "addr", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			var opErr *net.OpError
			if errors.As(err, &opErr) && errs.IsConnClosedErr(opErr.Err) {
				// Listener was closed
				return nil
			}
			return fmt.Errorf("error accepting new connection: %w", err)
		}

		go p.HandleConn(conn)
	}
}

func (p *GateGeyserPlugin) HandleConn(conn net.Conn) {
	// Handle the connection

	geyserConn := (&GeyserConnection{
		Conn: proxyproto.NewConn(conn),
		closeCb: func() {
			_ = conn.Close()
		},
	})

	p.mu.RLock()
	p.connections[geyserConn.RemoteAddr()] = geyserConn
	p.mu.RUnlock()

	go p.proxy.HandleConn(geyserConn)
}
