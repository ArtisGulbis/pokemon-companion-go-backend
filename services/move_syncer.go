package services

import (
	"sync"
	"time"
)

type MoveSyncer struct {
	client      MoveAPIClient
	repo        MoveRepo
	rateLimiter *time.Ticker
	syncedMoves map[int]bool // In-memory cache of synced move IDs
	mu          sync.Mutex   // Protects syncedMoves map
}

func NewMoveSyncer(client MoveAPIClient, repo MoveRepo, ticker *time.Ticker) *MoveSyncer {
	return &MoveSyncer{
		client:      client,
		repo:        repo,
		rateLimiter: ticker,
		syncedMoves: make(map[int]bool),
	}
}

func (m *MoveSyncer) SyncMove(id int) error {
	// Check cache first
	m.mu.Lock()
	if m.syncedMoves[id] {
		m.mu.Unlock()
		// Already synced in this session, skip API call
		return nil
	}
	m.mu.Unlock()

	// Not in cache, fetch from API
	move, err := m.client.FetchMove(id)
	if err != nil {
		return err
	}

	if err := m.repo.InsertMove(move); err != nil {
		return err
	}

	// Mark as synced in cache
	m.mu.Lock()
	m.syncedMoves[id] = true
	m.mu.Unlock()

	return nil
}
