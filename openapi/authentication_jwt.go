package openapi

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"os"
	"time"
)

func defaultSigningMethod() *jwt.SigningMethodHMAC {
	return jwt.SigningMethodHS256
}

func hmacSampleSecret() []byte {
	return []byte(os.Getenv("ACCESS_SECRET"))
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    time.Time
	RtExpires    time.Time
}

// Create a pair jwt tokens for authentication and refresh
func createToken(email string, name string) (string, string, error) {
	const accessDuration = time.Minute * 15
	const refreshDuration = time.Hour * 24 * 7

	// access access
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = uuid.NewV4().String()
	atClaims["user_email"] = email
	atClaims["user_name"] = name
	atClaims["exp"] = time.Now().UTC().Add(accessDuration).Unix()
	at := jwt.NewWithClaims(defaultSigningMethod(), atClaims)
	access, err := at.SignedString(hmacSampleSecret())
	if err != nil {
		return "", "", err
	}

	// refresh access
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = uuid.NewV4().String()
	atClaims["user_email"] = email
	atClaims["user_name"] = name
	rtClaims["exp"] = time.Now().UTC().Add(refreshDuration).Unix()
	rt := jwt.NewWithClaims(defaultSigningMethod(), rtClaims)
	refresh, err := rt.SignedString(hmacSampleSecret())
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
