package keydigital

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/byuoitav/connpool"
	"go.uber.org/zap"
)

var (
	ErrOutOfRange = errors.New("input or output is out of range")
	regGetInput   = regexp.MustCompile("Video Output  *: Input = ([0-9]{2}),")
)

// AudioVideoInputs .
func (vs *VideoSwitcher) AudioVideoInputs(ctx context.Context) (map[string]string, error) {
	vs.log.Info("Getting the current inputs")
	inputs := make(map[string]string)

	err := vs.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("unable to set connection deadline: %w", err)
		}

		cmd := []byte("STA\r\n")
		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return fmt.Errorf("unable to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("unable to write to connection: wrote %v/%v bytes", n, len(cmd))
		}

		var match [][]string
		for len(match) == 0 {
			buf, err := conn.ReadUntil(asciiCarriageReturn, deadline)
			if err != nil {
				return fmt.Errorf("unable to read from connection: %w", err)
			}

			match = regGetInput.FindAllStringSubmatch(string(buf), -1)
		}

		inputs[""] = strings.TrimPrefix(match[0][1], "0")
		return nil
	})
	if err != nil {
		return inputs, err
	}

	vs.log.Info("Got inputs", zap.Any("inputs", inputs))
	return inputs, nil
}

// SetAudioVideoInput .
func (vs *VideoSwitcher) SetAudioVideoInput(ctx context.Context, output, input string) error {
	output = "1"
	vs.log.Info("Setting audio video input", zap.String("output", output), zap.String("input", input))
	cmd := []byte(fmt.Sprintf("SPO0%sSI0%s\r\n", output, input))

	return vs.pool.Do(ctx, func(conn connpool.Conn) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			deadline = time.Now().Add(10 * time.Second)
		}

		if err := conn.SetDeadline(deadline); err != nil {
			return fmt.Errorf("unable to set connection deadline: %w", err)
		}

		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return fmt.Errorf("unable to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("unable to write to connection: wrote %v/%v bytes", n, len(cmd))
		}

		buf, err := conn.ReadUntil(asciiCarriageReturn, deadline)
		if err != nil {
			return fmt.Errorf("failed to read from connection: %w", err)
		}

		if strings.Contains(string(buf), "FAILED") {
			return ErrOutOfRange
		}

		vs.log.Info("Successfully set audio video input", zap.String("output", output), zap.String("input", input))
		return nil
	})
}
