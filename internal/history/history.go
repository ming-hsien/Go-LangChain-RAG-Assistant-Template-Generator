package history

import (
	"sync"

	"github.com/tmc/langchaingo/llms"
)

type HistoryManager struct {
	mu      sync.RWMutex
	history map[string][]llms.MessageContent
}

func NewHistoryManager() *HistoryManager {
	return &HistoryManager{
		history: make(map[string][]llms.MessageContent),
	}
}

func (m *HistoryManager) Get(sessionID string) []llms.MessageContent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if h, ok := m.history[sessionID]; ok {
		copyH := make([]llms.MessageContent, len(h))
		copy(copyH, h)
		return copyH
	}
	return nil
}

func (m *HistoryManager) Set(sessionID string, messages []llms.MessageContent) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.history[sessionID] = messages
}

func (m *HistoryManager) Clear(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.history, sessionID)
}
