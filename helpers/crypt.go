// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package helpers

import "golang.org/x/crypto/bcrypt"

func PasswordHash(pwd string) (string, error) {
	h, e := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if e != nil {
		return "", e
	}
	return string(h), nil
}

func PasswordVerify(hashedPwd string, pwd string) bool {
	e := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	if e != nil {
		return false
	}
	return true
}
