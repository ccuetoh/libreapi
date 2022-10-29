package rut

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ErrInvalidRUT = errors.New("invalid rut")
var ErrInvalidVD = errors.New("invalid verification figit")

type VD int8

const (
	VDNone VD = iota - 2
	VDK
	VD0
	VD1
	VD2
	VD3
	VD4
	VD5
	VD6
	VD7
	VD8
	VD9
)

type RUT struct {
	Digits []uint8
	VD     VD
}

var rutRegex = regexp.MustCompile("[0-9k]")

func parseRUT(rutStr string, ignoreVD bool) (RUT, error) {
	numsStr := rutRegex.FindAllString(strings.ToLower(rutStr), -1)
	if len(numsStr) < 7 || len(numsStr) > 10 {
		return RUT{VD: VDNone}, errors.Wrap(ErrInvalidRUT, "invalid length")
	}

	digitsStr := numsStr[:len(numsStr)-1]
	if ignoreVD {
		digitsStr = numsStr
	}

	var rut RUT
	for _, n := range digitsStr {
		num, err := strconv.ParseInt(n, 10, 8)
		if err != nil {
			return RUT{}, errors.Wrap(ErrInvalidRUT, fmt.Sprintf("invalid digit '%s'", n))
		}

		rut.Digits = append(rut.Digits, uint8(num))
	}

	if ignoreVD {
		rut.VD = VDNone
		return rut, nil
	}

	var err error
	rut.VD, err = parseVD(numsStr[len(numsStr)-1])
	if err != nil {
		return RUT{}, errors.Wrap(err, "invalid vd")
	}

	return rut, nil
}

func parseVD(vdStr string) (VD, error) {
	vdStr = strings.ToLower(vdStr)
	if vdStr == "k" {
		return VDK, nil
	}

	num, err := strconv.ParseInt(vdStr, 10, 8)
	if err != nil {
		return VDNone, errors.Wrap(ErrInvalidVD, fmt.Sprintf("invalid digit '%s'", vdStr))
	}

	if num < 0 || num > 10 {
		return VDNone, errors.Wrap(ErrInvalidVD, "vd should be k, or between 0 and 9")
	}

	return VD(num), nil
}

func generateRUT(min, max int) (RUT, error) {
	if min >= max {
		return RUT{}, errors.New("min should be lower than max")
	}

	digits := rand.Intn(max-min) + min
	digitsStr := strconv.Itoa(digits)

	rut, _ := parseRUT(digitsStr, true)
	rut.VD = rut.calculateVD()

	return rut, nil
}

func (r RUT) calculateVD() VD {
	seq := generateReverseSequence(len(r.Digits))
	var sum int
	for i, mask := range seq {
		sum += mask * int(r.Digits[i])
	}

	vd := 11 - (sum % 11)
	if vd == 10 {
		return VDK
	}

	if vd == 11 {
		return VD0
	}

	return VD(vd)
}

func (d VD) String() string {
	if d == VDNone {
		return "NONE"
	}

	if d == VDK {
		return "K"
	}

	return strconv.Itoa(int(d))
}

func (r RUT) IsValid() bool {
	expect := r.calculateVD()
	if expect == 11 {
		expect = 0
	}

	return r.VD == expect
}

func (r RUT) String() string {
	var builder strings.Builder
	for _, d := range r.Digits {
		builder.WriteString(strconv.Itoa(int(d)))
	}

	builder.WriteString("-")
	builder.WriteString(r.VD.String())

	return builder.String()
}

func generateReverseSequence(length int) (s []int) {
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
