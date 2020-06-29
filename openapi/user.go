package openapi

import (
	"encoding/base32"
)

type GamebaseUser struct {
	Name     string
	Email    string
	Password string
}

// construct a GamebaseUser from the data field of a v1 secret
func NewGamebaseUserFromSecretData(data map[string][]byte) GamebaseUser {
	name := string(data["name"])
	password := string(data["password"])

	email := string(data["email"])
	email = decodeEmail(email)
	return GamebaseUser{
		Name:     name,
		Email:    string(email),
		Password: password,
	}
}

func (user GamebaseUser) ToSecretData() map[string]string {
	return map[string]string{
		"name":     user.Name,
		"email":    encodeEmail(user.Email),
		"password": user.Password,
	}
}

func kubernetesFriendlyEncoding() *base32.Encoding {
	encoding := base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567")
	return encoding.WithPadding('_')
}

// encode the email as base32 with kubernetes friendly padding ('_' instead of '=').
func encodeEmail(email string) string {
	src := []byte(email)

	encoding := kubernetesFriendlyEncoding()
	buf := make([]byte, encoding.EncodedLen(len(src)))
	encoding.Encode(buf, src)

	return string(buf)
}

// decodeEmail the email as base32 with kubernetes friendly padding ('_' instead of '=').
func decodeEmail(email string) string {
	src := []byte(email)

	encoding := kubernetesFriendlyEncoding()
	buf := make([]byte, encoding.DecodedLen(len(src)))
	if i, err := encoding.Decode(buf, src); err != nil {
		panic("Could not decode email address \"" + email + "\" beginning with byte " + string(i))
	} else {
		return string(buf[:i])
	}
}
