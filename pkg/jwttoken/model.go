package jwttoken

import "errors"

var (
	NoValid = errors.New("NoValid")
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
}
