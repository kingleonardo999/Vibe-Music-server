package util

import "math/rand"

// GenRandomDigitalCode 生成n位纯数字随机验证码
func GenRandomDigitalCode(n int) string {
	digits := "0123456789"
	code := make([]byte, n)
	for i := 0; i < n; i++ {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}

// GenRandomCode 生成n位随机字符串验证码
func GenRandomCode(n int) string {
	letters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, n)
	for i := 0; i < n; i++ {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}
