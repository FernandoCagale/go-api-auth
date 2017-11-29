package models

import (
	shared "github.com/FernandoCagale/go-api-shared/src/validation"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) BeforeSave(scope *gorm.Scope) (err error) {
	password := []byte(u.Password)
	if pw, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost); err == nil {
		scope.SetColumn("password", pw)
		return nil
	}
	return err
}

func (u *User) ValidatePassword(password string) (valid bool) {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func (u *User) Validate() (errors map[string][]shared.Validation, ok bool) {
	errors = make(map[string][]shared.Validation)

	if u.Username == "" {
		errors["username"] = append(errors["username"], shared.Validation{
			Type:    "required",
			Message: "field is required",
		})
	}

	if len(u.Username) > 10 {
		errors["username"] = append(errors["username"], shared.Validation{
			Type:    "lenght-max",
			Message: "field lenght max 10",
		})
	}

	if len(u.Username) < 4 {
		errors["username"] = append(errors["username"], shared.Validation{
			Type:    "lenght-min",
			Message: "field lenght min 4",
		})
	}

	if u.Password == "" {
		errors["password"] = append(errors["password"], shared.Validation{
			Type:    "required",
			Message: "field is required",
		})
	}

	return errors, len(errors) == 0
}
