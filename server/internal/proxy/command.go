package proxy

import (
	"errors"
)

var (
	ErrUnknownProxyCmd = errors.New("unknown proxy command")
	ErrInvalidCommand  = errors.New("invalid command")
	ErrNilBuffer       = errors.New("buffer is nil")
	ErrDropped         = errors.New("data was dropped")
)

type ProxyCmdType byte

const (
	ProxyCmdUnknown ProxyCmdType = iota
	ProxyCmdStart
	ProxyCmdStop
	ProxyCmdStall
	ProxyCmdForward
	ProxyCmdDrop
)

var proxyCmdTypes = map[ProxyCmdType]string{
	ProxyCmdUnknown: "ProxyCmdUnknown",
	ProxyCmdStart:   "ProxyCmdStart",
	ProxyCmdStop:    "ProxyCmdStop",
	ProxyCmdStall:   "ProxyCmdStall",
	ProxyCmdForward: "ProxyCmdForward",
	ProxyCmdDrop:    "ProxyCmdDrop",
}

func (m ProxyCmdType) String() string {
	return proxyCmdTypes[m]
}

type ProxyCmd struct {
	Type ProxyCmdType `json:"type"`
	Data string       `json:"data"`
}

// Validate validates the data payloads of a ProxyCmd.
func (cmd *ProxyCmd) Validate() error {
	switch cmd.Type {
	case ProxyCmdStart, ProxyCmdStop, ProxyCmdDrop:
		if cmd.Data != "" {
			return ErrInvalidCommand
		}
	case ProxyCmdStall, ProxyCmdForward:
		if cmd.Data == "" {
			return ErrInvalidCommand
		}
	default:
		return ErrUnknownProxyCmd
	}

	return nil
}
