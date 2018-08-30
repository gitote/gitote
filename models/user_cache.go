package models

// MailResendCacheKey returns key used for cache mail resend.
func (u *User) MailResendCacheKey() string {
	return "MailResend_" + u.IDStr()
}

// TwoFactorCacheKey returns key used for cache two factor passcode.
// e.g. TwoFactor_1_012664
func (u *User) TwoFactorCacheKey(passcode string) string {
	return "TwoFactor_" + u.IDStr() + "_" + passcode
}
