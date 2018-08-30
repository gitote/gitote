// +build pam

package pam

import (
	"errors"

	"github.com/msteinert/pam"
)

func PAMAuth(serviceName, userName, passwd string) error {
	t, err := pam.StartFunc(serviceName, userName, func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff:
			return passwd, nil
		case pam.PromptEchoOn, pam.ErrorMsg, pam.TextInfo:
			return "", nil
		}
		return "", errors.New("Unrecognized PAM message style")
	})

	if err != nil {
		return err
	}

	if err = t.Authenticate(0); err != nil {
		return err
	}

	return nil
}
