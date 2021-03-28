package tcpscan

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsPortOpen(t *testing.T) {
	open := IsPortOpen("scanme.nmap.org", 22, 3 * time.Second)
	if !assert.True(t, open) {
		t.FailNow()
	}
	open = IsPortOpen("scanme.nmap.org", 8080, 3 * time.Second)
	if !assert.False(t, open) {
		t.FailNow()
	}
}

func TestGetOpenPorts(t *testing.T) {
	ports := []int {
		22,
		43,
		80,
		443,
		5432,
		9929,
	}
	open := GetOpenPorts("scanme.nmap.org", ports, 4)
	if !assert.Len(t, open, 3, "Incorrect number of returned open ports") {
		t.FailNow()
	}
	if !assert.Contains(t, open, 22, "SSH port not in open list") {
		t.FailNow()
	}
	if !assert.Contains(t, open, 80, "HTTP port not in open list") {
		t.FailNow()
	}
	if !assert.Contains(t, open, 9929, "nping-echo port not in open list") {
		t.FailNow()
	}
}