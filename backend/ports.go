package main

import (
	"sort"
	"strings"
	"sync"
	"time"

	gopsnet "github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type PortInfo struct {
	Port      uint32 `json:"port"`
	Protocol  string `json:"protocol"`  // TCP / UDP
	Status    string `json:"status"`    // 监听 / 已建立
	Family    string `json:"family"`    // IPv4 / IPv6
	Process   string `json:"process"`
	PID       int32  `json:"pid"`
	Address   string `json:"address"`
	StartedAt string `json:"startedAt"`
	HighRisk  bool   `json:"highRisk"`
}

type procMeta struct {
	name      string
	startedAt string
}

type Scanner struct {
	mu        sync.Mutex
	procCache map[int32]procMeta
}

func NewScanner() *Scanner {
	return &Scanner{procCache: make(map[int32]procMeta)}
}

func (s *Scanner) Scan(highRisk []int32) ([]PortInfo, error) {
	conns, err := gopsnet.Connections("all")
	if err != nil {
		return nil, err
	}

	risk := make(map[uint32]struct{}, len(highRisk))
	for _, p := range highRisk {
		if p > 0 {
			risk[uint32(p)] = struct{}{}
		}
	}

	out := make([]PortInfo, 0, len(conns))
	seen := make(map[string]struct{}, len(conns))

	for _, c := range conns {
		proto := protoName(c.Type, c.Family)
		if proto == "" {
			continue
		}
		// TCP: only listening sockets. UDP has no connection state — keep all.
		isUDP := strings.HasPrefix(proto, "UDP")
		if !isUDP && c.Status != "LISTEN" {
			continue
		}
		port := c.Laddr.Port
		if port == 0 {
			continue
		}

		family := "IPv4"
		if strings.Contains(c.Laddr.IP, ":") {
			family = "IPv6"
		}
		// Strip the v6 suffix from protoName so the column shows plain TCP/UDP per the design.
		proto = strings.TrimSuffix(proto, "6")

		key := proto + "|" + c.Laddr.IP + "|" + itoa(port) + "|" + itoa(uint32(c.Pid))
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}

		meta := s.lookupProc(c.Pid)

		status := "监听"
		if isUDP {
			status = "活动"
		}

		_, isRisk := risk[port]

		out = append(out, PortInfo{
			Port:      port,
			Protocol:  proto,
			Status:    status,
			Family:    family,
			Process:   meta.name,
			PID:       c.Pid,
			Address:   c.Laddr.IP,
			StartedAt: meta.startedAt,
			HighRisk:  isRisk,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Port != out[j].Port {
			return out[i].Port < out[j].Port
		}
		return out[i].Protocol < out[j].Protocol
	})
	return out, nil
}

func (s *Scanner) lookupProc(pid int32) procMeta {
	if pid <= 0 {
		return procMeta{name: "-", startedAt: "-"}
	}
	s.mu.Lock()
	if m, ok := s.procCache[pid]; ok {
		s.mu.Unlock()
		return m
	}
	s.mu.Unlock()

	m := procMeta{name: "-", startedAt: "-"}
	if p, err := process.NewProcess(pid); err == nil {
		if n, err := p.Name(); err == nil && n != "" {
			m.name = n
		}
		if ts, err := p.CreateTime(); err == nil && ts > 0 {
			m.startedAt = time.UnixMilli(ts).Format("2006-01-02 15:04:05")
		}
	}
	s.mu.Lock()
	s.procCache[pid] = m
	s.mu.Unlock()
	return m
}

// PrunePIDs drops cache entries for PIDs no longer present, called every minute or so.
func (s *Scanner) PrunePIDs(alive map[int32]struct{}) {
	s.mu.Lock()
	for pid := range s.procCache {
		if _, ok := alive[pid]; !ok {
			delete(s.procCache, pid)
		}
	}
	s.mu.Unlock()
}

func protoName(typ uint32, family uint32) string {
	// gopsutil uses syscall.SOCK_STREAM (1) for TCP, SOCK_DGRAM (2) for UDP
	switch typ {
	case 1:
		if family == 10 || family == 30 {
			return "TCP6"
		}
		return "TCP"
	case 2:
		if family == 10 || family == 30 {
			return "UDP6"
		}
		return "UDP"
	}
	return ""
}

func itoa(u uint32) string {
	if u == 0 {
		return "0"
	}
	var b [10]byte
	i := len(b)
	for u > 0 {
		i--
		b[i] = byte('0' + u%10)
		u /= 10
	}
	return string(b[i:])
}
