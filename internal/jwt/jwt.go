package jwt

import "github.com/dgrijalva/jwt-go"

// Claims type that uses for JSON decoding.
type Claims struct {
	UID   string `json:"uid"`
	Chips int64  `json:"chips"`
	Bet   int64  `json:"bet"`
}

// Valid validates Claims. There's no validation rules, always valid.
func (c Claims) Valid() error { return nil }

// Service an interface for dealing with jwt tokens.
type Service interface {
	Encode(claims Claims) (string, error)
	Decode(tokenString string) (Claims, error)
}

type basicService struct {
	secret []byte
}

// NewService returns Service implementation.
func NewService(secret []byte) Service {
	return &basicService{secret: secret}
}

// Encode encodes claims and returns jwt token.
func (svc *basicService) Encode(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(svc.secret)
}

// Decode decodes jwt token and returns its claims.
func (svc *basicService) Decode(tokenString string) (Claims, error) {
	c := Claims{}
	if _, err := jwt.ParseWithClaims(tokenString, &c, svc.keyFunc); err != nil {
		return c, err
	}

	return c, nil
}

func (svc *basicService) keyFunc(token *jwt.Token) (interface{}, error) {
	return svc.secret, nil
}
