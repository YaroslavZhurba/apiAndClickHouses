package repository

import (
	"sync"
)

const ipV4size = 32

type IpMaskRepository struct {
	clock           Clock              // clockInterface for Testing
	maskSize        uint32             // 255.255.255.34/24, MaskSize is 24
	requestLimit    uint32             // limit number of requests
	cooldownSeconds float64            // cooldown in seconds
	mu              sync.Mutex         // mutex need when we insert in map
	requestsCount   map[uint32]*Subnet // key is a Subnet mask, contains count of requests & timeOfFirstRequest
}

// if mask ip 1111111000.0000 then will get 1111111
// in other worlds remove zeros from mask
func (r *IpMaskRepository) getMask(ip uint32) uint32 {
	return ip >> (ipV4size - r.maskSize)
}

func NewIpRepository(maskSize, requestLimit uint32, cooldownSeconds float64) *IpMaskRepository {
	return &IpMaskRepository{
		clock:           RealClock{},
		requestsCount:   make(map[uint32]*Subnet),
		maskSize:        maskSize,
		requestLimit:    requestLimit,
		cooldownSeconds: cooldownSeconds,
	}
}

// function check a limit for a subnet
func (r *IpMaskRepository) ProcessRequest(ip uint32) bool {
	now := r.clock.Now()
	mask := r.getMask(ip)
	subnet := r.getSubnet(mask)

	subnet.mu.Lock() // take lock only for a subnet (do not lock other subnets)
	defer subnet.mu.Unlock()

	if subnet.cooldown { //check if status cooldown
		timeDiff := now.Sub(subnet.firstRequestTime)
		if timeDiff.Seconds() < r.cooldownSeconds { // check if cooldown time is passed
			return false
		}
		subnet.cooldown = false // if cooldown time passed refressh data
		subnet.requestsNumber = 1
		subnet.firstRequestTime = now
		return true
	}

	if subnet.requestsNumber == 0 { // check if subnet a new
		subnet.requestsNumber = 1
		subnet.firstRequestTime = now
		subnet.cooldown = false
		return true
	}
	subnet.requestsNumber++

	if subnet.requestsNumber >= r.requestLimit { // check our limit
		timeDiff := now.Sub(subnet.firstRequestTime)
		if timeDiff.Seconds() > r.cooldownSeconds {
			subnet.requestsNumber = 1
			subnet.firstRequestTime = now
			return true
		}
		subnet.cooldown = true
		subnet.firstRequestTime = now
		return false
	}

	return true
}

// function clear limit counter
func (r *IpMaskRepository) ProcessClear(ip uint32) {
	now := r.clock.Now()
	mask := r.getMask(ip)
	subnet := r.getSubnet(mask)

	subnet.mu.Lock() // take lock only for a subnet (do not lock other subnets)
	defer subnet.mu.Unlock()

	subnet.requestsNumber = 0
	subnet.firstRequestTime = now
	subnet.cooldown = false
}

func (r *IpMaskRepository) getSubnet(mask uint32) *Subnet {
	subnet, ok := r.requestsCount[mask] // check is exists
	if !ok {
		r.mu.Lock() // if not exists need to take lock, otherwise data ccould be overwritten if someone else would be faster
		defer r.mu.Unlock()
		_, exists := r.requestsCount[mask] // need to check, beause first check was unlocked
		if !exists {
			r.requestsCount[mask] = &Subnet{}
		}
		subnet = r.requestsCount[mask]

	}
	return subnet
}
