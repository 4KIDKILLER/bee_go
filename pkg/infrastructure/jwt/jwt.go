package jwt

import (
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

var secret = "07311726_BEE"

type BeeClaims struct {
	Username string `json:"username"`
	UserId   int    `json:"userId"`
	jwtv5.RegisteredClaims
}

type BeeJwt struct{}

func NewBeeJwt() (beeJwt *BeeJwt) {
	return &BeeJwt{}
}

func (beeJwt *BeeJwt) GenerateToken(username string, userId int) (token string, err error) {
	beeClaims := BeeClaims{
		Username: username,
		UserId:   userId,
		RegisteredClaims: jwtv5.RegisteredClaims{
			//超时时间，从现在开始往后一小时
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Hour)),
			//签发时间
			IssuedAt: jwtv5.NewNumericDate(time.Now()),
			//发行人
			Issuer: "BEE_V1CTOR",
		},
	}
	token, err = jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, beeClaims).SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return
}

func (beeJwt *BeeJwt) ParseToken(tokenString string) (*BeeClaims, error) {
	token, err := jwtv5.ParseWithClaims(tokenString, &BeeClaims{}, func(token *jwtv5.Token) (any, error) {
		return []byte(secret), nil
	}, jwtv5.WithValidMethods([]string{"HS265"}))

	if token != nil {
		if claims, ok := token.Claims.(*BeeClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, err
}
