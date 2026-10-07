package service

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func ValidateNIM(nim string) error {
	matched, _ := regexp.MatchString(`^\d{12}$`, nim)

	if !matched {
		return fmt.Errorf("NIM harus terdiri dari 12 digit")
	}

	return nil
}

func ValidateAngkatan(angkatan int) error {
	currentYear := time.Now().Year()

	if angkatan < 1000 || angkatan > currentYear {
		return fmt.Errorf("angkatan harus 4 digit dan tidak boleh melebihi tahun berjalan")
	}

	return nil
}

func ValidateIPK(ipk float64) error {
	if ipk < 0 || ipk > 4 {
		return fmt.Errorf("IPK harus berada pada rentang 0 sampai 4")
	}

	return nil
}

func ValidatePositiveID(id string) (int64, error) {
	value, err := strconv.ParseInt(id, 10, 64)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("ID tidak valid")
	}

	return value, nil
}