package keydigital

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/byuoitav/connpool"
)

type info struct {
	IPAddress       string
	MACAddress      string
	FirmwareVersion string
}

var (
	regIPAddr  = regexp.MustCompile("Host IP Address = ([0-9]{3}.[0-9]{3}.[0-9]{3}.[0-9]{3})")
	regMacAddr = regexp.MustCompile("MAC Address = ([A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2})")
	regVersion = regexp.MustCompile("Version : ([0-9]+.[0-9]+)")
)

func (vs *VideoSwitcher) Info(ctx context.Context) (interface{}, error) {
	vs.log.Info("Getting info")
	var info info

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

		for len(info.FirmwareVersion) == 0 || len(info.IPAddress) == 0 || len(info.MACAddress) == 0 {
			buf, err := conn.ReadUntil(asciiCarriageReturn, deadline)
			if err != nil {
				return fmt.Errorf("unable to read from connection: %w", err)
			}

			if len(info.FirmwareVersion) == 0 {
				match := regVersion.FindAllStringSubmatch(string(buf), -1)
				if len(match) > 0 {
					info.FirmwareVersion = match[0][1]
				}
			}

			if len(info.IPAddress) == 0 {
				match := regIPAddr.FindAllStringSubmatch(string(buf), -1)
				if len(match) > 0 {
					info.IPAddress = match[0][1]
				}
			}

			if len(info.MACAddress) == 0 {
				match := regMacAddr.FindAllStringSubmatch(string(buf), -1)
				if len(match) > 0 {
					info.MACAddress = match[0][1]
				}
			}
		}

		return nil
	})

	return info, err
}

func (vs *VideoSwitcher) Healthy(ctx context.Context) error {
	_, err := vs.AudioVideoInputs(ctx)
	return err
}
