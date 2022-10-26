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
	for i, num := range nums {
		switch num {
		case "k":
			if i != len(nums)-1 {
				return RUT{}, ErrInvalidRUT
			}

			rut2 = append(rut2, 10)
			continue
		default:
			numInt, err := strconv.Atoi(num)
			if err != nil {
				return RUT{}, ErrInvalidRUT
			}

			rut2 = append(rut2, numInt)
		}
	}

	return rut2, nil
}

func (r RUT) CalculateVD(ignoreLast bool) int {
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

	return 11 - (sum % 11)
}

func VDToString(vd int) string {
	switch vd {
	case 10:
		return "K"
	case 11:
		return "0"
	default:
		return strconv.Itoa(vd)
	}
}

func (r RUT) GetVDString() string {
	val := strconv.Itoa(r[len(r)-1])
	switch val {
	case "10":
		return "K"
	case "11":
		return "0"
	default:
		return val
	}
}

func (r RUT) IsValid() bool {
	expect := r.CalculateVD(true)
	if expect == 11 { // -0 RUTs
		expect = 0
	}

	return r[len(r)-1] == expect
}

func (r RUT) String() string {
	if len(r) < 7 || len(r) > 10 {
		return ""
	}

	var digits string
	for _, d := range r[:len(r)-1] {
		digits += strconv.Itoa(d)
	}

	return digits + "-" + VDToString(r.CalculateVD(true))
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

	return s + "-" + VDToString(r[len(r)-1])
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
