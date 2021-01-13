package keydigital

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestGetInput(t *testing.T) {
	is := is.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// vs := NewVideoSwitcher("BRMB-382-SW1.byu.edu", WithLogger(zap.NewExample()))
	vs := NewVideoSwitcher("JFSB-B132-SW1.byu.edu", WithLogger(zap.NewExample()))

	inputs, err := vs.AudioVideoInputs(ctx)
	is.NoErr(err)
	t.Logf("inputs: %v\n", inputs)
}
