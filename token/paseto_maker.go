package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func NewPasetoMaker(symmetricKey []byte) (Maker,error){
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil,fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		symmetricKey: symmetricKey,
		paseto: paseto.NewV2(),
	},nil
}

func (q *PasetoMaker) CreateTokenPaseto(username string, duration time.Duration) (string,*Payload,error) {
	payload := NewPayload(username,duration)

	token,err := q.paseto.Encrypt(q.symmetricKey,payload,nil)

	return token,payload,err
}

func (q *PasetoMaker) VerifyTokenPaseto(token string)(*Payload,error) {
	payload := &Payload{}

	err := q.paseto.Decrypt(token,q.symmetricKey,payload,nil)

	if err != nil {
		return nil,errInvalid
	}

	if payload.valid() {
		return nil,errExpired
	}

	return payload,nil
}