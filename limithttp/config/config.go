package config

import (
	"flag"
)

type Options struct {
	MaskSize        uint32
	RequestLimit    uint32
	CooldownSeconds float64
}

func (o *Options) Read() error {
	var (
		maskSize     uint
		requestLimit uint
	)
	flag.UintVar(&maskSize, "mask-size", 24, "number of bits in a mask")
	flag.UintVar(&requestLimit, "limit", 100, "limit number of requests")
	flag.Float64Var(&o.CooldownSeconds, "cooldown", 60.0, "cooldown in seconds")
	flag.Parse()

	o.MaskSize = uint32(maskSize)
	o.RequestLimit = uint32(requestLimit)
	return nil
}
