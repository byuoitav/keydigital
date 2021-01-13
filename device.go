package keydigital

import (
	"context"
	"net"
	"time"

	"github.com/byuoitav/connpool"
	"go.uber.org/zap"
)

const (
	asciiCarriageReturn = 0x0d
)

type VideoSwitcher struct {
	pool *connpool.Pool
	log  *zap.Logger
}

func NewVideoSwitcher(addr string, opts ...Option) *VideoSwitcher {
	options := &options{
		ttl:   30 * time.Second,
		delay: 250 * time.Millisecond,
		log:   zap.NewNop(),
	}

	for _, o := range opts {
		o.apply(options)
	}

	return &VideoSwitcher{
		log: options.log,
		pool: &connpool.Pool{
			TTL:   options.ttl,
			Delay: options.delay,
			NewConnection: func(ctx context.Context) (net.Conn, error) {
				dial := net.Dialer{}
				return dial.DialContext(ctx, "tcp", addr+":23")
			},
			Logger: options.log.Sugar(),
		},
	}
}
