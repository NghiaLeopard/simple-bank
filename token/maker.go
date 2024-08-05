package token

import "time"

type Maker interface {
	CreateTokenPaseto(username string, duration time.Duration) (string,*Payload,error)

	VerifyTokenPaseto(token string) (*Payload,error)
}