package states

import (
	"context"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
)

type Repository interface {
	CreateState(context.Context, *entities.State) error
	GetState(context.Context, string, int) (*entities.State, error)
	DeleteState(context.Context, string, int) error
}

type State struct {
	rep   Repository
	sType int
}

func New(rep Repository, sType int) *State {
	return &State{rep: rep, sType: sType}
}

func (s State) Add(id string, data interface{}) error {
	state := &entities.State{
		ID:   id,
		Type: s.sType,
		Data: data,
	}

	return s.rep.CreateState(context.Background(), state)
}

func (s State) Get(id string) (interface{}, error) {
	state, err := s.rep.GetState(context.Background(), id, s.sType)
	if err != nil {
		return nil, err
	}

	return state.Data, nil
}

func (s State) Delete(id string) error {
	return s.rep.DeleteState(context.Background(), id, s.sType)
}
