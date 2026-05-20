package messaging

import (
	"sync"

	"github.com/pocketbase/pocketbase/models"
)

// Conversation states for the onboarding and posting flows.
const (
	stateIdle            = ""
	stateSignupName      = "signup_name"
	stateSignupEmail     = "signup_email"
	stateSignupPassword  = "signup_password"
	stateSignupTower     = "signup_tower"
	stateSignupApartment  = "signup_apartment"
	stateSignupCreateApt  = "signup_create_apt"
	stateSignupPasskey    = "signup_passkey"
	stateLoginEmail      = "login_email"
	stateLoginPassword   = "login_password"
	statePostType        = "post_type"
	statePostTitle       = "post_title"
)

type session struct {
	User    *models.Record
	State   string
	Pending map[string]string // holds in-progress data across state transitions
}

func newSession() *session {
	return &session{Pending: make(map[string]string)}
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

func getOrCreate(from string) *session {
	sessions.Lock()
	defer sessions.Unlock()
	if s, ok := sessions.m[from]; ok {
		return s
	}
	s := newSession()
	sessions.m[from] = s
	return s
}

func clearSession(from string) {
	sessions.Lock()
	defer sessions.Unlock()
	sessions.m[from] = newSession()
}
