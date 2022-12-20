package tokens

import (
	"log"
	"time"
	"os"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = os.Getenv("JWTSECRET")

type JWTClaim struct {
	Login    string;
	Role	string;
	jwt.StandardClaims
}

func GenerateJWT(login string, role string, expirationTime time.Time) string {
	claims:= &JWTClaim{
		Login: login,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		log.Fatalf("couldn't handle this token: %#v", err)
	}

	return tokenString
}

func ValidateToken(signedToken string) (bool, string, string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	claims, ok := token.Claims.(*JWTClaim)

	if err != nil {
		return false, claims.Login, claims.Role
	}

	if !ok {
		log.Fatalf("couldn't parse claims")
	}
	
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return false, claims.Login, claims.Role
	}
	return true, claims.Login, claims.Role
}

func AccessToken(login string, role string) string {
	accessTokenDuration := time.Now().Add(1 * time.Minute)
	return GenerateJWT(login, role, accessTokenDuration)
}

func RefreshToken(login string, role string) string {
	refreshTokenDuration := time.Now().Add(1 * time.Hour)
	return GenerateJWT(login, role, refreshTokenDuration)
}