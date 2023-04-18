package repository

import (
	"testing"
	"time"
)

type MockClock struct {
	t time.Time
}

func NewMocClock() Clock {
	return &MockClock{t: time.Unix(1681737372, 0).UTC()}
}

func (m *MockClock) Now() time.Time {
	return m.t
}

func (m *MockClock) AddTime(seconds int) {
	duration := time.Duration(seconds) * time.Second
	m.t = m.t.Add(duration)
}

func getIp(a, b, c, d int) uint32 {
	return uint32(a*256*256*256 + b*256*256 + c*256 + d)
}

func Test_ProcessRequests_TooManyRequests(t *testing.T) {
	repo := IpMaskRepository{
		clock:           NewMocClock(),
		requestsCount:   make(map[uint32]*Subnet),
		maskSize:        24,
		requestLimit:    5,
		cooldownSeconds: 10,
	}
	ip1 := getIp(254, 255, 44, 13)
	ip2 := getIp(254, 255, 44, 65)
	repo.ProcessRequest(ip1)
	repo.ProcessRequest(ip1)
	repo.ProcessRequest(ip1)
	repo.ProcessRequest(ip2)
	if repo.ProcessRequest(ip2) != false {
		t.Errorf("want false, got true")
	}
}

func Test_ProcessRequests_TooManyRequests_ThenWait(t *testing.T) {
	repo := IpMaskRepository{
		clock:           NewMocClock(),
		requestsCount:   make(map[uint32]*Subnet),
		maskSize:        24,
		requestLimit:    5,
		cooldownSeconds: 10,
	}
	ip1 := getIp(254, 255, 44, 13)
	ip2 := getIp(254, 255, 44, 65)
	repo.ProcessRequest(ip1)
	repo.ProcessRequest(ip1)
	repo.ProcessRequest(ip1)
	repo.ProcessRequest(ip2)
	repo.ProcessRequest(ip2)
	repo.clock.AddTime(15)
	if repo.ProcessRequest(ip1) != true {
		t.Errorf("want true, got false")
	}
}

func Test_ProcessRequests_ManyRequestWithLongTime(t *testing.T) {
	repo := IpMaskRepository{
		clock:           NewMocClock(),
		requestsCount:   make(map[uint32]*Subnet),
		maskSize:        24,
		requestLimit:    5,
		cooldownSeconds: 10,
	}
	ip1 := getIp(254, 255, 44, 13)
	ip2 := getIp(254, 255, 44, 65)
	repo.ProcessRequest(ip1)
	repo.clock.AddTime(3)
	repo.ProcessRequest(ip1)
	repo.clock.AddTime(3)
	repo.ProcessRequest(ip1)
	repo.clock.AddTime(3)
	repo.ProcessRequest(ip2)
	repo.clock.AddTime(3)
	repo.ProcessRequest(ip2)

	if repo.ProcessRequest(ip1) != true {
		t.Errorf("want true, got false")
	}
}

func Test_ProcessRequests_FirstRequest(t *testing.T) {
	repo := IpMaskRepository{
		clock:           NewMocClock(),
		requestsCount:   make(map[uint32]*Subnet),
		maskSize:        24,
		requestLimit:    5,
		cooldownSeconds: 10,
	}
	ip1 := getIp(254, 255, 44, 13)
	if repo.ProcessRequest(ip1) != true {
		t.Errorf("want true, got false")
	}
}
