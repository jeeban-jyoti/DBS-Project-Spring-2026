package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

var store = struct {
	m map[string]string
	sync.RWMutex
}{
	m: make(map[string]string),
}

func GenerateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func Create(sessionID, email string) {
	store.Lock()
	defer store.Unlock()
	store.m[sessionID] = email
}

func Get(sessionID string) (string, bool) {
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
