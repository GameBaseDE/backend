package openapi

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"log"
	"os"
	"time"
)

func defaultSigningMethod() *jwt.SigningMethodHMAC {
	return jwt.SigningMethodHS256
}

func hmacSampleSecret() []byte {
	//FIXME only call once on startup
	key := make([]byte, 32)
	secret := os.Getenv("ACCESS_SECRET")
	if secret == "" {
		_, err := rand.Read(key)
		if err != nil {
			log.Fatal("Could not create random hmacSampleSecret")
		}
		os.Setenv("ACCESS_SECRET", base64.StdEncoding.EncodeToString(key))
		return key
	} else {
		readKey, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			return []byte(secret)
		}
		return readKey
	}
}

type userClaims struct {
	TokenUuid    string `json:"token_uuid,omitempty"`
	UserEmail    string `json:"user_email,omitempty"`
	UserName     string `json:"user_name,omitempty"`
	UserGravatar string `json:"user_gravatar,omitempty"`
	jwt.StandardClaims
}

// Create a pair jwt tokens for authentication and refresh
func createToken(user GamebaseUser) (string, string, error) {
	const accessDuration = time.Minute * 15
	const refreshDuration = time.Hour * 24 * 7

	now := time.Now().UTC()

	// access access
	var err error
	atClaims := userClaims{
		TokenUuid:    uuid.NewV4().String(),
		UserEmail:    user.Email,
		UserName:     user.Name,
		UserGravatar: user.Gravatar,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(accessDuration).Unix(),
		},
	}
	at := jwt.NewWithClaims(defaultSigningMethod(), atClaims)
	access, err := at.SignedString(hmacSampleSecret())
	if err != nil {
		return "", "", err
	}

	// refresh access
	rtClaims := userClaims{
		TokenUuid: uuid.NewV4().String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(refreshDuration).Unix(),
		},
	}
	rt := jwt.NewWithClaims(defaultSigningMethod(), rtClaims)
	refresh, err := rt.SignedString(hmacSampleSecret())
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
