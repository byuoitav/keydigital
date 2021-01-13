package keydigital

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestGetInfo(t *testing.T) {
	is := is.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// vs := NewVideoSwitcher("BRMB-382-SW1.byu.edu", WithLogger(zap.NewExample()))
	vs := NewVideoSwitcher("JFSB-B132-SW1.byu.edu", WithLogger(zap.NewExample()))

	info, err := vs.Info(ctx)
	is.NoErr(err)
	t.Logf("info: %+v\n", info)
}

func TestHealthy(t *testing.T) {
	is := is.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vs := NewVideoSwitcher("BRMB-382-SW1.byu.edu", WithLogger(zap.NewExample()))
	// vs := NewVideoSwitcher("JFSB-B132-SW1.byu.edu", WithLogger(zap.NewExample()))

	is.NoErr(vs.Healthy(ctx))
}
