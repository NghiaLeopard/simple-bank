package token

import (
	"errors"
	"time"
)

type Payload struct {
	Username string `json:"username"`
	IssuedAt time.Time `json:"issuedAt"`
	Expired time.Time `json:"expired"`
}

var (
	errInvalid = errors.New("token is invalid")
	errExpired = errors.New("token is expired")
)

func NewPayload(username string,duration time.Duration) *Payload {
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	return &Payload{
		Username: username,
		IssuedAt: issuedAt,
		Expired: expiredAt,
	}
}

func (p *Payload) valid() bool {
	return time.Now().After(p.Expired)
}