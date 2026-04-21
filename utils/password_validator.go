package utils

import (
	"errors"
	"regexp"
)

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Le mot de passe doit contenir au minimum 8 caractères.")
	}

	var (
		upper   = regexp.MustCompile(`[A-Z]`)
		lower   = regexp.MustCompile(`[a-z]`)
		number  = regexp.MustCompile(`[0-9]`)
		special = regexp.MustCompile(`[!@#%\$%\^&\*\.]`)
	)

	if !upper.MatchString(password) {
		return errors.New("Le mot de passe doit contenir au minimum une majuscule.")
	}

	if !lower.MatchString(password) {
		return errors.New("Le mot de passe doit contenir au minimum une miniscule.")
	}

	if !number.MatchString(password) {
		return errors.New("Le mot de passe doit contenir au minimum un chiffre.")
	}

	if !special.MatchString(password) {
		return errors.New("Le mot de passe doit contenir au minimum un caractère spécial.")
	}

	return nil // Si toutes les règles sont validées.

}
