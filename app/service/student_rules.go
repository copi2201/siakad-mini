package service

func GetSKSLimit(ipk float64) int {
	if ipk >= 3.00 {
		return 24
	}

	if ipk >= 2.50 {
		return 21
	}

	return 18
}