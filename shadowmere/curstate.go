package shadowmere

import (
	"sync"
)

type set map[string]bool

/**
 * Non-persistent urrent services state - which users are identified,
 * what nicks are connected, and so forth.
 */
type CurState struct {
	glock sync.RWMutex  // Eew

	// Nicks that have identified to NickServ
	identifiedNicks set

	// New nicks that have connected, and are registered, but
	// have yet to identify to NickServ
	newNicks set
}

func NewCurState() *CurState {
	return &CurState{
		identifiedNicks: make(set),
		newNicks:        make(set),
	}
}

func (cs *CurState) Identify(nick string) {
	cs.glock.Lock()
	defer cs.glock.Unlock()
	cs.identifiedNicks[nick] = true
	delete(cs.newNicks, nick)
}
func (cs *CurState) Unidentify(nick string) {
	cs.glock.Lock()
	defer cs.glock.Unlock()
	delete(cs.identifiedNicks, nick)
}
func (cs *CurState) NewNick(nick string) {
	cs.glock.Lock()
	defer cs.glock.Unlock()
	cs.newNicks[nick] = true
}

func (cs *CurState) IsIdentified(nick string) bool {
	cs.glock.RLock()
	defer cs.glock.RUnlock()
	return cs.identifiedNicks[nick] == true
}
func (cs *CurState) IsNew(nick string) bool {
	cs.glock.RLock()
	defer cs.glock.RUnlock()
	return cs.newNicks[nick] == true
}