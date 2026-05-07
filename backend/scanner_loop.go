package main

import (
	"encoding/json"
	"log"
	"sync/atomic"
	"time"
)

type ScannerLoop struct {
	scanner *Scanner
	store   *ConfigStore
	hub     *Hub
	reload  chan struct{}
	stop    chan struct{}
	cur     atomic.Int64 // current interval in ms (for read-only display)
}

func NewScannerLoop(s *Scanner, cs *ConfigStore, h *Hub) *ScannerLoop {
	l := &ScannerLoop{
		scanner: s,
		store:   cs,
		hub:     h,
		reload:  make(chan struct{}, 1),
		stop:    make(chan struct{}),
	}
	l.cur.Store(int64(cs.Snapshot().IntervalMs))
	return l
}

func (l *ScannerLoop) Reload() {
	select {
	case l.reload <- struct{}{}:
	default:
	}
}

func (l *ScannerLoop) Stop() {
	close(l.stop)
}

func (l *ScannerLoop) Run() {
	pruneTicker := time.NewTicker(60 * time.Second)
	defer pruneTicker.Stop()

	for {
		cfg := l.store.Snapshot()
		l.cur.Store(int64(cfg.IntervalMs))
		l.tick(cfg)

		ticker := time.NewTicker(time.Duration(cfg.IntervalMs) * time.Millisecond)
	inner:
		for {
			select {
			case <-l.stop:
				ticker.Stop()
				return
			case <-l.reload:
				ticker.Stop()
				break inner
			case <-pruneTicker.C:
				// noop here; pruning handled in tick via PrunePIDs would need alive set
			case <-ticker.C:
				cfg := l.store.Snapshot()
				l.cur.Store(int64(cfg.IntervalMs))
				l.tick(cfg)
			}
		}
	}
}

type wsFrame struct {
	Type      string     `json:"type"`
	Ports     []PortInfo `json:"ports"`
	Timestamp int64      `json:"timestamp"`
}

func (l *ScannerLoop) tick(cfg Config) {
	ports, err := l.scanner.Scan(cfg.HighRiskPorts)
	if err != nil {
		log.Printf("scan error: %v", err)
		return
	}
	frame := wsFrame{Type: "ports", Ports: ports, Timestamp: time.Now().UnixMilli()}
	payload, err := json.Marshal(frame)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}
	l.hub.Broadcast(payload)

	alive := make(map[int32]struct{}, len(ports))
	for _, p := range ports {
		alive[p.PID] = struct{}{}
	}
	l.scanner.PrunePIDs(alive)
}
