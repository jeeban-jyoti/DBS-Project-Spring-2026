package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type SessionData struct {
	Email string
	Role  string
}

var store = struct {
	m map[string]SessionData
	sync.RWMutex
}{
	m: make(map[string]SessionData),
}

func GenerateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func Create(sessionID, email, role string) {
	store.Lock()
	defer store.Unlock()
	store.m[sessionID] = SessionData{
		Email: email,
		Role:  role,
	}
}

func Get(sessionID string) (SessionData, bool) {
	store.RLock()
	defer store.RUnlock()
	val, ok := store.m[sessionID]
	return val, ok
}

func Delete(sessionID string) {
	store.Lock()
	defer store.Unlock()
	delete(store.m, sessionID)
}
