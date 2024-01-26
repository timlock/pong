package repository

import (
	"pong/entity"
	"sync"

	"github.com/google/uuid"
)

type UserStore struct {
	mu      sync.RWMutex
	players map[uuid.UUID]*entity.Player
}

func NewUserStore() UserStore {
	return UserStore{mu: sync.RWMutex{}, players: make(map[uuid.UUID]*entity.Player)}
}

func (p *UserStore) Add(player *entity.Player) {
	p.mu.Lock()
	p.players[player.Id] = player
	p.mu.Unlock()
}

func (p *UserStore) Get(id uuid.UUID) (*entity.Player, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	player, ok := p.players[id]
	return player, ok
}
