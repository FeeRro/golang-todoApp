package domain

import (
	"fmt"
	"regexp"

	core_error "github.com/FeeRro/golang-todoApp/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

func NewUser(id int, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		phoneNumber,
	)
}

func (u User) Validate() error {
	fullNameLenght := len([]rune(u.FullName))
	if fullNameLenght < 3 || fullNameLenght > 100 {
		return fmt.Errorf("invalid full name lenght %d: %w", fullNameLenght, core_error.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberlenght := len([]rune(*u.PhoneNumber))
		if phoneNumberlenght < 10 || phoneNumberlenght > 15 {
			return fmt.Errorf("invalid phone number lenght %d: %w", phoneNumberlenght, core_error.ErrInvalidArgument)
		}

		re := regexp.MustCompile(`^\+?[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid phone number format: %w", core_error.ErrInvalidArgument)
		}
	}
	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"'FullName' can't be patched to NULL: %w",
			core_error.ErrInvalidArgument,
		)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
