package repository

import (
	"sync"
	"time"
)

type Subnet struct {
	mu               sync.Mutex // used for concurrency
	requestsNumber   uint32     // type to store number of requests
	cooldown         bool       // cooldown state, true = Too Manu Requests
	firstRequestTime time.Time
}
