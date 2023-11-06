package services

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	key []byte
}

func NewJWTService(key []byte) *JwtService {
	return &JwtService{
		key: key,
	}
}
func (j JwtService) CreateKey(email, role string) (*string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": email,
		"role": role,
	})
	signed, err := t.SignedString(j.key)
	if err != nil {
		return nil, err
	}
	return &signed, nil
}
