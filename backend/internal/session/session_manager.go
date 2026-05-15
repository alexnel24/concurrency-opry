package session

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

const defaultSessionTTLHours = 24

type Session struct {
	ID         string
	LastActive time.Time
}

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	ttl      time.Duration
	wg       sync.WaitGroup
}

func NewSessionManager() *SessionManager {
	ttlHours := defaultSessionTTLHours
	if envVal := os.Getenv("SESSION_TTL_HOURS"); envVal != "" {
		if parsed, err := strconv.Atoi(envVal); err == nil && parsed > 0 {
			ttlHours = parsed
		}
	}
	return &SessionManager{
		sessions: make(map[string]*Session),
		ttl:      time.Duration(ttlHours) * time.Hour,
	}
}

func (sm *SessionManager) Create() string {
	id := generateID()

	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.sessions[id] = &Session{ID: id, LastActive: time.Now()}
	return id
}

func (sm *SessionManager) IsValid(id string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	s, ok := sm.sessions[id]
	if !ok {
		return false
	}
	if time.Since(s.LastActive) > sm.ttl {
		delete(sm.sessions, id)
		return false
	}
	s.LastActive = time.Now()
	return true
}

func (sm *SessionManager) Delete(id string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, ok := sm.sessions[id]; !ok {
		return false
	}
	delete(sm.sessions, id)
	return true
}

func (sm *SessionManager) StartBackgroundSessionCleanup(ctx context.Context) {
	sm.wg.Add(1)
	go func() {
		defer sm.wg.Done()
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				sm.mu.Lock()
				for id, s := range sm.sessions {
					if time.Since(s.LastActive) > sm.ttl {
						delete(sm.sessions, id)
					}
				}
				sm.mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
