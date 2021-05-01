package api

import "github.com/pkg/errors"

func Api() (string, error) {
	return "hello from api", errors.New("not implemented")
}
