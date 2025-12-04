package util

import "golang.org/x/crypto/bcrypt"

// EncryptPassword 加密密码
func EncryptPassword(password string) (string, error) {
	encryptedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encryptedBytes), err
}

// ComparePassword 对比明文密码和加密后的密码是否匹配
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
