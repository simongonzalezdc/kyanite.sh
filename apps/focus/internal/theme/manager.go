package theme

import "sync"

var (
	globalManager *Manager
	once          sync.Once
)

type Manager struct {
	mu      sync.RWMutex
	current Theme
}

func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{current: Default()}
	})
	return globalManager
}

func (m *Manager) SetTheme(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.current = GetTheme(id)
}

func (m *Manager) Current() Theme {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

func (m *Manager) Next() Theme {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := ListThemes()
	cur := m.current
	curID := ""
	for id, t := range Registry {
		if t.Name == cur.Name {
			curID = id
			break
		}
	}
	for i, id := range ids {
		if id == curID {
			next := ids[(i+1)%len(ids)]
			m.current = Registry[next]
			return m.current
		}
	}
	m.current = Default()
	return m.current
}

func (m *Manager) Previous() Theme {
	m.mu.Lock()
	defer m.mu.Unlock()
	ids := ListThemes()
	cur := m.current
	curID := ""
	for id, t := range Registry {
		if t.Name == cur.Name {
			curID = id
			break
		}
	}
	for i, id := range ids {
		if id == curID {
			prev := ids[(i-1+len(ids))%len(ids)]
			m.current = Registry[prev]
			return m.current
		}
	}
	m.current = Default()
	return m.current
}
