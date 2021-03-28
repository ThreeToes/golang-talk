package tcpscan

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// scanPool convenience wrapper to
type scanPool struct {
	dialLock    *sync.Mutex
	portsLock   *sync.Mutex
	done        *sync.WaitGroup
	portsToDial []int
	address     string
	openPorts   []int
	timeout     time.Duration
}

// START DIALWORKER OMIT
// dialWorker processes the sp.portsToDial list and adds the open ports to sp.openPorts OMIT
func (sp *scanPool) dialWorker() {
	defer sp.done.Done() // OMIT
	var open bool        // OMIT
	var portToDial int   // OMIT
	// ...
	for {
		sp.dialLock.Lock()
		if sp.portsToDial == nil {
			sp.dialLock.Unlock()
			return
		}
		portToDial = sp.portsToDial[0]
		if len(sp.portsToDial) == 1 {
			sp.portsToDial = nil
		}else {
			sp.portsToDial = sp.portsToDial[1:]
		}
		sp.dialLock.Unlock()
		open = IsPortOpen(sp.address, portToDial, sp.timeout)
		if open {
			sp.portsLock.Lock()
			sp.openPorts = append(sp.openPorts, portToDial)
			sp.portsLock.Unlock()
		}
	}
}
// END DIALWORKER OMIT

// IsPortOpen scan a port on a remote IP and see if it's openPorts
func IsPortOpen(addr string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", addr, port), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// GetOpenPorts on addr from the list ports
func GetOpenPorts(addr string, ports []int, taskCount int) []int {
	sp := &scanPool{
		dialLock:    &sync.Mutex{},
		portsLock:   &sync.Mutex{},
		done:        &sync.WaitGroup{},
		openPorts:   nil,
		address:     addr,
		portsToDial: ports,
	}
	sp.done.Add(taskCount)
	for i := 0; i < taskCount; i++ {
		go sp.dialWorker()
	}
	sp.done.Wait()
	return sp.openPorts
}