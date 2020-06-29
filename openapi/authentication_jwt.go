package openapi

import (
	"github.com/denisbrodbeck/machineid"
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
	if secret := os.Getenv("ACCESS_SECRET"); secret != "" {
		return []byte(secret)
	} else {
		log.Println("Could not read ENV $ACCESS_SECRET")
		id, err := machineid.ProtectedID("GameBase")
		if err != nil {
			log.Fatal("Failed to alternatively read unique OS id")
		}
		return []byte(id)
	}
}

type Claims struct {
	TokenUuid string `json:"token_uuid,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	jwt.StandardClaims
}

// Create a pair jwt tokens for authentication and refresh
func createToken(email string, name string) (string, string, error) {
	const accessDuration = time.Minute * 15
	const refreshDuration = time.Hour * 24 * 7

	now := time.Now().UTC()

	// access access
	var err error
	atClaims := Claims{
		TokenUuid: uuid.NewV4().String(),
		UserEmail: email,
		UserName:  name,
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
	rtClaims := Claims{
		TokenUuid: uuid.NewV4().String(),
		UserEmail: email,
		UserName:  name,
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
