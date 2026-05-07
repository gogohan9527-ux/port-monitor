package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	cfg     *ConfigStore
	auth    *Auth
	scanner *Scanner
	hub     *Hub
	loop    *ScannerLoop
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type loginResp struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.Username != hardcodedUser || req.Password != hardcodedPass {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}
	tok, exp, err := s.auth.Issue(req.Username)
	if err != nil {
		http.Error(w, "issue token failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, loginResp{Token: tok, ExpiresAt: exp.UnixMilli()})
}

func (s *Server) handlePorts(w http.ResponseWriter, r *http.Request) {
	cfg := s.cfg.Snapshot()
	ports, err := s.scanner.Scan(cfg.HighRiskPorts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{
		"ports":     ports,
		"timestamp": time.Now().UnixMilli(),
	})
}

func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, s.cfg.Snapshot())
}

type updateSettingsReq struct {
	IntervalMs    int     `json:"intervalMs"`
	HighRiskPorts []int32 `json:"highRiskPorts"`
}

func (s *Server) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req updateSettingsReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if req.IntervalMs > 0 && req.IntervalMs < 200 {
		http.Error(w, "intervalMs 不能小于 200", http.StatusBadRequest)
		return
	}
	cfg, err := s.cfg.Update(req.IntervalMs, req.HighRiskPorts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.loop.Reload()
	writeJSON(w, cfg)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 4)}
	s.hub.Register(client)

	go func() {
		defer func() {
			s.hub.Unregister(client)
			conn.Close()
		}()
		conn.SetReadLimit(1024)
		conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(120 * time.Second))
			return nil
		})
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	go func() {
		ping := time.NewTicker(30 * time.Second)
		defer ping.Stop()
		for {
			select {
			case msg, ok := <-client.send:
				if !ok {
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					return
				}
			case <-ping.C:
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}
