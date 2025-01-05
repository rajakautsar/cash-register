package auth

import (
    "time"
    "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

func GenerateAdminToken(username, role string) (string, error) {
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: username,
        Role:     role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}