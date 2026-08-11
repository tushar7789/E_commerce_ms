package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func GetTokens(userID, accessSecret string, accessTTL time.Duration, refreshSecret string, refreshTTL time.Duration) (TokenPair, error) {
	accessToken, err := CreateTokens(userID, accessSecret, accessTTL)
	if err != nil {
		log.Print(err)
		return TokenPair{}, err
	}

	refreshToken, err := CreateTokens(userID, refreshSecret, refreshTTL)
	if err != nil {
		log.Print(err)
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func CreateTokens(userID, secret string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var durationRe = regexp.MustCompile(`(?i)^([0-9]+(?:\.[0-9]+)?)(ns|us|µs|ms|s|m|h|d|w)$`)

func mustDuration(v string) time.Duration {
	d, err := parseDurationDW(v)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}

func parseDurationDW(s string) (time.Duration, error) {
	m := durationRe.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("invalid duration: %q", s)
	}

	val, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, err
	}

	switch m[2] {
	case "ns":
		return time.Duration(val), nil
	case "us", "µs":
		return time.Duration(val * float64(time.Microsecond)), nil
	case "ms":
		return time.Duration(val * float64(time.Millisecond)), nil
	case "s":
		return time.Duration(val * float64(time.Second)), nil
	case "m":
		return time.Duration(val * float64(time.Minute)), nil
	case "h":
		return time.Duration(val * float64(time.Hour)), nil
	case "d":
		return time.Duration(val * float64(24*time.Hour)), nil
	case "w":
		return time.Duration(val * float64(7*24*time.Hour)), nil
	default:
		return 0, fmt.Errorf("invalid duration unit: %q", m[2])
	}
}
