package messaging

import (
	"sync"

	"github.com/pocketbase/pocketbase/models"
)

// session holds in-memory auth state per "from" address (phone number / simulator ID).
// This is intentionally simple — fine for local dev and single-instance production.
// For multi-instance deployments, replace with a Redis-backed store.
type session struct {
	User *models.Record
}

var sessions = struct {
	sync.RWMutex
	m map[string]*session
}{m: make(map[string]*session)}

func getSession(from string) *session {
	sessions.RLock()
	defer sessions.RUnlock()
	return sessions.m[from]
}

func setSession(from string, s *session) {
	sessions.Lock()
	defer sessions.Unlock()
	sessions.m[from] = s
}

func clearSession(from string) {
	sessions.Lock()
	defer sessions.Unlock()
	delete(sessions.m, from)
}
