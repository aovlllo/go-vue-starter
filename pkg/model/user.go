package model

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// User is the structure of a user
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`

	Name       string `json:"name"`
	SecondName string `json:"secondName"`
	Birth      string `json:"birth"`
	Sex        string `json:"sex"`
	Interests  string `json:"interests"`
	City       string `json:"city"`
}

// HashPassword hashed the password of the user
func (u *User) HashPassword() error {
	key, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(key)

	return nil
}

// MatchPassword returns true if the hashed user password matches the password
func (u *User) MatchPassword(password string) bool {
	key, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	keyStr := string(key)
	_ = keyStr

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err == nil {
		return true
	}

	return false
}

var sexValues = []string{"male", "female", "non binary"}

func (u *User) VerifyFields() error {
	switch {
	case u.Email == "":
		return fmt.Errorf("email address is missing")
	case u.Name == "":
		return fmt.Errorf("name is missing")
	}
	var validSex bool
	for _, s := range sexValues {
		if u.Sex == s {
			validSex = true
			break
		}
	}
	if !validSex {
		return fmt.Errorf("sex is not valid")
	}
	return nil
}
