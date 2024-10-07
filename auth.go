package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stonoy/my_remainder/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func hashFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

func compareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type jwtCusromClaims struct {
	jwt.RegisteredClaims
}

func generateToken(user database.User, jwtSecret string) (string, error) {
	claims := jwtCusromClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "my_remainder",
			ID:        fmt.Sprintf("%v", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSecret))
	return ss, err
}

func validateToken(tokenString, jwtSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCusromClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwtCusromClaims); ok && token.Valid {
		return claims.ID, nil
	} else {
		return "", fmt.Errorf("error parsing token")
	}
}

func getTokenFromHeader(r *http.Request) string {
	authHaeder := r.Header.Get("authorization")

	if len(authHaeder) < 1 {
		return ""
	}

	authSlice := strings.Fields(authHaeder)
	if authSlice[0] == "Bearer" && len(authSlice) == 2 {
		return authSlice[1]
	} else {
		return ""
	}
}

type authFunc func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) authTokenToUser(theFunc authFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token
		authToken := getTokenFromHeader(r)
		if authToken == "" {
			replyWithError("no token provided", 401, w)
			return
		}

		userIDStr, err := validateToken(authToken, cfg.jwt_secret)
		if err != nil {
			replyWithError("no valid token provided", 401, w)
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			replyWithError("not a valid user", 401, w)
			return
		}

		theUser, err := cfg.dbQ.GetUserByID(r.Context(), userID)
		if err != nil {
			if err == sql.ErrNoRows {
				replyWithError("No such user found", 400, w)
			} else {
				replyWithError("error in GetUserByID", 500, w)
			}
			return
		}

		theFunc(w, r, theUser)

	}
}
