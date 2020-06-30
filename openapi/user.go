package openapi

import (
	"encoding/base32"
)

type GamebaseUser struct {
	Name     string
	Email    string
	Password string
	Gravatar string
}

// construct a GamebaseUser from the data field of a v1 secret
func NewGamebaseUserFromSecretData(email string, data map[string][]byte) GamebaseUser {
	name := string(data["name"])
	password := string(data["password"])
	gravatar := string(data["gravatar"])

	return GamebaseUser{
		Name:     name,
		Email:    email,
		Password: password,
		Gravatar: gravatar,
	}
}

func (user GamebaseUser) ToSecretData() map[string]string {
	return map[string]string{
		"name":     user.Name,
		"password": user.Password,
		"gravatar": user.Gravatar,
	}
}

func kubernetesFriendlyEncoding() *base32.Encoding {
	encoding := base32.NewEncoding("abcdefghijklmnopqrstuvwxyz123456")
	return encoding.WithPadding('0')
}

// encode the email as base32 with kubernetes friendly padding ('0' instead of '=').
func encodeEmail(email string) string {
	src := []byte(email)

	encoding := kubernetesFriendlyEncoding()
	buf := make([]byte, encoding.EncodedLen(len(src)))
	encoding.Encode(buf, src)

	return string(buf)
}

// decodeEmail the email as base32 with kubernetes friendly padding ('0' instead of '=').
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
