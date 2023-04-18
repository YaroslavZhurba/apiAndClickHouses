package server

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

func getIp(ctx *fasthttp.RequestCtx) (uint32, error) {
	xForwardedFor := string(ctx.Request.Header.Peek("X-Forwarded-For")) // get X-Forwarded-For
	ipString := strings.Split(xForwardedFor, ",")[0]                    // get ip as a string
	if ipString == "" {
		return 0, fmt.Errorf("empty X-Forwarded-For")
	}
	ipInt, err := ipStringToUInt32(ipString)
	if err != nil {
		return 0, err
	}
	return ipInt, nil
}
