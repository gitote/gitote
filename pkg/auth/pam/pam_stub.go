// +build !pam

package pam

import (
	"errors"
)

func PAMAuth(serviceName, userName, passwd string) error {
	return errors.New("PAM not supported")
}
