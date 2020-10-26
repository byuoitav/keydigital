package keydigital

import (
	"context"
	"fmt"
	"regexp"

	"github.com/byuoitav/connpool"
)

var (
	regIPAddr  = regexp.MustCompile("Host IP Address = ([0-9]{3}.[0-9]{3}.[0-9]{3}.[0-9]{3})")
	regMacAddr = regexp.MustCompile("MAC Address = ([A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2}:[A-Z,0-9]{2})")
	regVersion = regexp.MustCompile("Version : ([0-9]+.[0-9]+)")
)

type Info struct {
	IPAddress       string
	MACAddress      string
	FirmwareVersion string
}

//Info .
func (vs *VideoSwitcher) Info(ctx context.Context) (Info, error) {
	var resp Info

	if vs.Pool.Logger != nil {
		vs.Pool.Logger.Infof("getting hardware info")
	}

	err := vs.Pool.Do(ctx, func(conn connpool.Conn) error {
		cmd := []byte("STA\r\n")
		n, err := conn.Write(cmd)
		switch {
		case err != nil:
			return fmt.Errorf("failed to write to connection: %w", err)
		case n != len(cmd):
			return fmt.Errorf("failed to write to connection: wrote %v/%v bytes", n, len(cmd))
		}

		var match [][]string
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("no deadline set")
		}
		for len(match) == 0 {
			buf, err := conn.ReadUntil(carriageReturn, deadline)
			if err != nil {
				return fmt.Errorf("failed to read from connection: %w", err)
			}

			// TODO make sure match[0] exists (and match[0][1])

			// Mac Address
			match = regMacAddr.FindAllStringSubmatch(string(buf), -1)
			if len(match) >= 1 {
				resp.MACAddress = match[0][1]
			}

			// Version
			match = regVersion.FindAllStringSubmatch(string(buf), -1)
			if len(match) >= 1 {
				resp.FirmwareVersion = match[0][1]
			}

			// IP Address
			match = regIPAddr.FindAllStringSubmatch(string(buf), -1)
			if len(match) >= 1 {
				resp.IPAddress = match[0][1]
			}
		}

		return nil
	})

	if err != nil {
		return resp, err
	}

	return resp, nil
}
