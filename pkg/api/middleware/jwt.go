package middleware

import (
	"carbide-registry-api/pkg/api/utils"
	"carbide-registry-api/pkg/objects"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	// "github.com/rs/zerolog/log"
)

func Authenticate(w http.ResponseWriter, user objects.User) error {
	err := setAuthCookie(w, user)
	return err
}

func Authorized(w http.ResponseWriter, r *http.Request) bool {
	_, err := verifyJWT(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		utils.RespondWithJSON(w, "user is unauthorized")
		return false
	}
	return true
}

func setAuthCookie(w http.ResponseWriter, user objects.User) error {
	token, err := generateJWT(user)
	if err != nil {
		return err
	}
	ck := http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &ck)
	return nil
}

// generate JWT for given user
func generateJWT(user objects.User) (string, error) {
	secret := os.Getenv("JWTSECRET")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	claims["userid"] = fmt.Sprint(user.Id)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyJWT(r *http.Request) (int64, error) {
	secret := os.Getenv("JWTSECRET")
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return 0, err
		}
		return 0, err
	}
	tokenString := c.Value
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid JWT")
		}
		return []byte(secret), nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("failed to parse JWT claims")
	}
	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return 0, errors.New("token expired")
	}
	userIdString := claims["userid"].(string)
	userID, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return 0, errors.New("failed to parse userid")
	}
	return userID, nil
}

func terminateJWT() {
	// replace jwt with another that expires immediately
}
