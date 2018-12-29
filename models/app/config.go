package app

import (
	trinary "github.com/iotaledger/iota.go/trinary"
)

// Config - Config
type Config struct {
	Endpoint string         `json:"endpoint"`
	Seed     trinary.Trytes `json:"seed"`
	MWM      uint64         `json:"mwm"`
	Depth    uint64         `json:"depth"`
}
