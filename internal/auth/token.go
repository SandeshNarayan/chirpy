package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)



func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){


	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
        Subject: userID.String(),
	})

	tokenString, err := claims.SignedString([]byte(tokenSecret))
	if err!=nil{
        return "", err
    }
	return tokenString, nil
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	
	
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC);!ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(tokenSecret), nil
    })

    if err!= nil ||!token.Valid {
        return uuid.Nil, err
    }

    claims, ok := token.Claims.(*jwt.RegisteredClaims)
    if!ok || claims.ExpiresAt ==nil || claims.ExpiresAt.Time.Before(time.Now()) {
        return uuid.Nil, fmt.Errorf("token is expired or invalid")
    }

    return uuid.Parse(claims.Subject)
} 


func GetBearerToken(headers http.Header) (string, error){
	
	bearer:= headers.Get("Authorization")
	if bearer==""{
		return "", fmt.Errorf("authorization header is missing")
	}
	
	parts := strings.SplitN(bearer, " ", 2)
	if len(parts)!=2 || parts[0] !="Bearer"{
		return "", fmt.Errorf("invalid bearer token format")
	}
	return parts[1], nil
}