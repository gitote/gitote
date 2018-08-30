package errors

import "fmt"

type TwoFactorNotFound struct {
	UserID int64
}

func IsTwoFactorNotFound(err error) bool {
	_, ok := err.(TwoFactorNotFound)
	return ok
}

func (err TwoFactorNotFound) Error() string {
	return fmt.Sprintf("two-factor authentication does not found [user_id: %d]", err.UserID)
}

type TwoFactorRecoveryCodeNotFound struct {
	Code string
}

func IsTwoFactorRecoveryCodeNotFound(err error) bool {
	_, ok := err.(TwoFactorRecoveryCodeNotFound)
	return ok
}

func (err TwoFactorRecoveryCodeNotFound) Error() string {
	return fmt.Sprintf("two-factor recovery code does not found [code: %s]", err.Code)
}
