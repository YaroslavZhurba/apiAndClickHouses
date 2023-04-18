package server

import (
	"strconv"
	"strings"
)

// transform ip string to uint32
// sample "0.0.10.9" = 2569
func ipStringToUInt32(ipString string) (uint32, error) {
	ip := uint32(0)
	strs := strings.Split(ipString, ".")
	for _, s := range strs {
		num, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		ip *= 256
		ip += uint32(num)
	}
	return ip, nil
}
