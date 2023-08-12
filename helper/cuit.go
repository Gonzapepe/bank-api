package helper

import (
	"errors"
	"fmt"
	"strconv"
)

// Cuit creates the Cuit of the person based on their national Identity and gender
func Cuit(dni int64, gender int) (int64, error) {
	dniString := strconv.Itoa(int(dni))

	var AB string
	// 0 for male, 1 for female
	switch gender {
	case 0:
		AB = "20"
	case 1:
		AB = "27"
	default:
		return 0, errors.New("invalid gender")
	}

	AB0, _ := strconv.Atoi(string(AB[0]))
	AB1, _ := strconv.Atoi(string(AB[1]))
	calc := AB0*5 + AB1*4

	multipliers := []int{3, 2, 7, 6, 5, 4, 3, 2}
	for i := 0; i < 8; i++ {
		digit, _ := strconv.Atoi(string(dniString[i]))
		calc += digit * multipliers[i]
	}

	Z := 11 - (calc % 11)

	if Z == 10 {
		return Cuit(dni, gender)
	}

	if Z == 11 {
		Z = 0
	}

	result := fmt.Sprintf("%s%s%d", AB, dniString, Z)
	res, err := strconv.ParseInt(result, 10, 64)

	if err != nil {
		return 0, err
	}
	return res, nil

}