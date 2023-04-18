package main

import (
	"limithttp/config"
	"limithttp/src/server"
	"log"
)

func main() {
	opt := config.Options{}
	if err := opt.Read(); err != nil {
		log.Fatal(err)
	}
	server.NewServer(opt.MaskSize, opt.RequestLimit, opt.CooldownSeconds).Run(8080)

}
