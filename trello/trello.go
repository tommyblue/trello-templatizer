package trello

import "errors"

type Trello struct {
	key   string
	token string
}

func New(key, token string) (*Trello, error) {
	if key == "" || token == "" {
		return nil, errors.New("key and token can't be empty")
	}

	return &Trello{
		key:   key,
		token: token,
	}, nil
}
