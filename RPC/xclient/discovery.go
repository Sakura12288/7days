package xclient

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelectMode SelectMode = iota
	RoundRobinMode
)

type Discovery interface {
	Refresh() error
	Get(SelectMode) (string, error)
	GetAll() ([]string, error)
	Update([]string) error
}

type MultiServersDiscovery struct {
	r       *rand.Rand
	mu      sync.RWMutex
	index   int
	servers []string
}

var _ Discovery = (*MultiServersDiscovery)(nil)

func NewMultiServerDiscovery(servers []string) *MultiServersDiscovery {
	d := &MultiServersDiscovery{
		servers: servers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	d.index = d.r.Intn(math.MaxInt32 - 1)
	return d
}

func (d *MultiServersDiscovery) Refresh() error {
	return nil
}

func (d *MultiServersDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.RUnlock()
	d.servers = servers
	return nil
}

func (d *MultiServersDiscovery) Get(mode SelectMode) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.servers) == 0 {
		return "", errors.New("没有服务")
	}
	switch mode {
	case RandomSelectMode:
		index := d.r.Intn(len(d.servers))
		return d.servers[index], nil
	case RoundRobinMode:
		server := d.servers[d.index%len(d.servers)]
		d.index = (d.index + 1) % len(d.servers)
		return server, nil
	default:
		return "", fmt.Errorf("该模式不支持 %d", mode)
	}
}

func (d *MultiServersDiscovery) GetAll() ([]string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}
