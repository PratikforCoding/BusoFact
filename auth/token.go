package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrorNoAuthHeaderIncluded = errors.New("not auth header included in request")

func MakeAccessToken(id string, jwtSecret string, expiresIn time.Duration) (string, error) {
    signingKey := []byte(jwtSecret)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
        Subject:   id,
        Issuer:    "busofact-access", // Match the issuer here
    })

    return token.SignedString(signingKey)
}

func MakeRefreshToken(id string, jwtSecret string, expiresIn time.Duration) (string, error) {
    signingKey := []byte(jwtSecret)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
        Subject:   id,
        Issuer:    "busofact-refresh", // Match the issuer here
    })

    return token.SignedString(signingKey)
}

func GetTokenFromCookie(r *http.Request, tokenSecret string) (string, error) {
    // Get access_token from cookie
    cookie, err := r.Cookie("access_token")
    if err != nil {
		if err == http.ErrNoCookie {
			cookie, err = r.Cookie("refresh_token")
			if err != nil {
				if err == http.ErrNoCookie {
					return "", errors.New("no refreshcookie included in request")
				}
			}
			newAccessToken, err := RefreshToken(cookie.Value, tokenSecret)
			if err != nil {
				return "", err
			}
			return newAccessToken, nil
		}
		
		return "", errors.New("couldn't get cookie from request")
	}

    return cookie.Value, nil
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {return []byte(tokenSecret), nil},
	)
	if err != nil {
		return "", err
	}

	issuer, err :=  token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}

	if issuer != "busofact-access" {
		return "", errors.New("token is refresh token")
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	return userIDString, nil
} 

func RefreshToken(tokenString, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != "busofact-refresh" {
		return "", errors.New("invalid issuer")
	}

	if err != nil {
		return "", err
	}

	newToken, err := MakeAccessToken(
		userIDString,
		tokenSecret,
		time.Hour,
	)
	if err != nil {
		return "", err
	}

	return newToken, nil
}