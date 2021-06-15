package rut

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidRUT = errors.New("invalid rut")

type RUT []int

func ParseRUT(rut string) (RUT, error) {
	re := regexp.MustCompile("[0-9k]")

	nums := re.FindAllString(strings.ToLower(rut), -1)
	if len(nums) < 7 || len(nums) > 10 {
		return RUT{}, ErrInvalidRUT
	}

	var rut2 RUT
	for _, num := range nums {
		if num == "k" {
			rut2 = append(rut2, 0)
			continue
		}

		numInt, err := strconv.Atoi(num)
		if err != nil {
			return RUT{}, ErrInvalidRUT
		}

		rut2 = append(rut2, numInt)
	}

	return rut2, nil
}

func (r RUT) CalculateValidationDigit(ignoreLast bool) int {
	var digits []int
	if ignoreLast {
		digits = r[:len(r)-1]
	} else {
		digits = r
	}

	seq := GetReverseSequence(len(digits))
	var sum int
	for i, mask := range seq {
		sum += mask * r[i]
	}

	res := 11 - (sum % 11)
	if res == 10 {
		return 0 // K
	}

	return res
}

func (r RUT) IsValid() bool {
	if len(r) < 7 || len(r) > 10 {
		return false
	}

	return r[len(r)-1] == r.CalculateValidationDigit(true)
}

func (r RUT) String() string {
	if len(r) < 7 || len(r) > 10 {
		return ""
	}

	var digits string
	for _, d := range r[:len(r)-1] {
		digits += strconv.Itoa(d)
	}

	return digits + "-" + r.GetValidationDigit()
}

func (r RUT) PrettyString() string {
	if len(r) < 7 || len(r) > 10 {
		return ""
	}

	digits := r[:len(r)-1]

	var digitsStr []string
	for _, digit := range digits {
		digitsStr = append(digitsStr, strconv.Itoa(digit))
	}

	var s string
	for i := len(digits) % 3; i <= len(digits); i += 3 {
		x := i - 3
		if x < 0 {
			x = 0
		}

		s += strings.Join(digitsStr[x:i], "")
		if i+3 <= len(digits) && len(s) > 0 {
			s += "."
		}
	}

	return s + "-" + r.GetValidationDigit()
}

func GetReverseSequence(length int) (s []int) {
	for len(s) < length {
		if len(s) == 0 {
			s = append(s, 2)
			continue
		}

		if s[len(s)-1] >= 7 {
			s = append(s, 2)
			continue
		}

		s = append(s, s[len(s)-1]+1)
	}

	for i := len(s)/2 - 1; i >= 0; i-- {
		opp := len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}

	return s
}

func (r RUT) GetValidationDigit() string {
	val := strconv.Itoa(r[len(r)-1])
	if val == "0" {
		val = "K"
	}

	return val
}
