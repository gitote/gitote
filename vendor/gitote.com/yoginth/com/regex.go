package com

import "regexp"

const (
	regex_email_pattern        = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
	regex_strict_email_pattern = `(?i)[A-Z0-9!#$%&'*+/=?^_{|}~-]+` +
		`(?:\.[A-Z0-9!#$%&'*+/=?^_{|}~-]+)*` +
		`@(?:[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?\.)+` +
		`[A-Z0-9](?:[A-Z0-9-]*[A-Z0-9])?`
	regex_url_pattern = `(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`
)

var (
	regex_email        *regexp.Regexp
	regex_strict_email *regexp.Regexp
	regex_url          *regexp.Regexp
)

func init() {
	regex_email = regexp.MustCompile(regex_email_pattern)
	regex_strict_email = regexp.MustCompile(regex_strict_email_pattern)
	regex_url = regexp.MustCompile(regex_url_pattern)
}

// validate string is an email address, if not return false
// basically validation can match 99% cases
func IsEmail(email string) bool {
	return regex_email.MatchString(email)
}

// validate string is an email address, if not return false
// this validation omits RFC 2822
func IsEmailRFC(email string) bool {
	return regex_strict_email.MatchString(email)
}

// validate string is a url link, if not return false
// simple validation can match 99% cases
func IsUrl(url string) bool {
	return regex_url.MatchString(url)
}
