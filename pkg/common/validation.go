package common

import (
	"regexp"
)

// E-posta validasyon fonksiyonu
func IsEmailValid(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Telefon numarası validasyon fonksiyonu (sadece rakam ve 10-15 karakter arası)
func IsPhoneValid(phone string) bool {
	re := regexp.MustCompile(`^[0-9]{10,15}$`)
	return re.MatchString(phone)
}
