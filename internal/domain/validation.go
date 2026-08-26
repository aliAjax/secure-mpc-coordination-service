package domain

import (
	"net/url"
	"regexp"
)

var idPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.:-]{2,127}$`)

func ValidateID(v string) error {
	if !idPattern.MatchString(v) {
		return ErrInvalid
	}
	return nil
}
func ValidateTenant(v string) error {
	if len(v) < 2 || len(v) > 64 {
		return ErrInvalid
	}
	return ValidateID(v)
}
func ValidateProtocol(name, version string) error {
	if err := ValidateID(name); err != nil {
		return err
	}
	if err := ValidateID(version); err != nil {
		return err
	}
	return nil
}
func ValidateCallback(raw string) error {
	u, e := url.Parse(raw)
	if e != nil || u.Scheme != "https" || u.Host == "" {
		return ErrInvalid
	}
	return nil
}
func IsTerminal(s ComputationStatus) bool {
	return s == StatusSucceeded || s == StatusAborted || s == StatusExpired
}
