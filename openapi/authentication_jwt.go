package openapi

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"os"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    time.Time
	RtExpires    time.Time
}

// check if the access token is expired
func (t *TokenDetails) isExpired() bool {
	return time.Now().UTC().After(t.AtExpires)
}

// check if the refresh token is expired
func (t *TokenDetails) isExpiredForever() bool {
	return time.Now().UTC().After(t.RtExpires)
}

var loggedInUsers = make(map[string]*TokenDetails)

func createToken(email string) (*TokenDetails, error) {
	td := &TokenDetails{
		AtExpires:   time.Now().UTC().Add(time.Minute * 15),
		AccessUuid:  uuid.NewV4().String(),
		RtExpires:   time.Now().UTC().Add(time.Hour * 24 * 7),
		RefreshUuid: uuid.NewV4().String(),
	}

	// access token
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_email"] = email
	atClaims["exp"] = td.AtExpires.Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_email"] = email
	rtClaims["exp"] = td.RtExpires.Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	loggedInUsers[email] = td

	return td, nil
}
