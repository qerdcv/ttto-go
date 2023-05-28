package domain

import (
	"crypto/sha256"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/pbkdf2"
)

type User struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Username, validation.Required, validation.Length(3, 30)),
		validation.Field(&u.Password, validation.Required, validation.Length(4, 30)),
	)
}

func (u *User) HashPassword(secret []byte) {
	u.Password = fmt.Sprintf("%x", string(pbkdf2.Key([]byte(u.Password), secret, 100_000, 32, sha256.New)))
}
