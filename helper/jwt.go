package helper

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("SomeSecret")

type CustomClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(id string) (string, error) {
	mySigningKey = []byte(mySigningKey)

	claims := CustomClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 1000000).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(mySigningKey)

	return ss, err
}

func VerifyToken(t string) (string, error) {
	token, err := jwt.ParseWithClaims(t, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.ID, nil
	}

	return "", nil
}
