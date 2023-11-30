package entities

import (
	"encoding/json"
	"fmt"
)

type State struct {
	ID   string `db:"id"`
	Type int    `db:"s_type"`
	Data any    `db:"data"`
}

func (s *State) EncodeData() error {
	bs, err := json.Marshal(&s.Data)
	if err != nil {
		return err
	}

	s.Data = bs
	return nil
}

func (s *State) DecodeData() error {
	bs, ok := s.Data.(string)
	if !ok {
		return fmt.Errorf("error parsing data")
	}

	var data any
	if err := json.Unmarshal([]byte(bs), &data); err != nil {
		return err
	}

	s.Data = data
	return nil
}
