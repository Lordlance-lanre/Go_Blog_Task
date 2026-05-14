package utils
import (
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

var SecretKey = "your_secret_key"

func GenerateJWT(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Issuer: issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	return claims.SignedString([]byte(SecretKey))	
}

func ParseJWT(cookie string) (string, error) {
 token, error  := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
 if error != nil || !token.Valid {
	 return "", error
 }
 return token.Claims.(*jwt.StandardClaims).Issuer, nil
}