package rut

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRut(t *testing.T) {
	rut, err := parseRUT("123123-1", false)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{1, 2, 3, 1, 2, 3}, rut.Digits)
	assert.Equal(t, VD1, rut.VD)
	assert.Equal(t, "123123-1", rut.String())
	assert.False(t, rut.IsValid())

	rut, err = parseRUT("123.123.123-1", false)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{1, 2, 3, 1, 2, 3, 1, 2, 3}, rut.Digits)
	assert.Equal(t, VD1, rut.VD)
	assert.Equal(t, "123123123-1", rut.String())
	assert.False(t, rut.IsValid())

	rut, err = parseRUT("123.123.123-1", true)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{1, 2, 3, 1, 2, 3, 1, 2, 3, 1}, rut.Digits)
	assert.Equal(t, rut.VD, VDNone)
	assert.False(t, rut.IsValid())

	rut, err = parseRUT("asdasd", true)
	assert.Error(t, err)
	assert.Equal(t, ([]uint8)(nil), rut.Digits)
	assert.Equal(t, VDNone, rut.VD)
	assert.False(t, rut.IsValid())

	rut, err = parseRUT("123123-a", false)
	assert.Error(t, err)
	assert.Equal(t, ([]uint8)(nil), rut.Digits)
	assert.Equal(t, VDNone, rut.VD)
	assert.False(t, rut.IsValid())

	rut, err = parseRUT("123123-0", false)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{1, 2, 3, 1, 2, 3}, rut.Digits)
	assert.Equal(t, VD0, rut.VD)
	assert.Equal(t, "123123-0", rut.String())
	assert.False(t, rut.IsValid())

	rut, err = parseRUT("123123-K", false)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{1, 2, 3, 1, 2, 3}, rut.Digits)
	assert.Equal(t, VDK, rut.VD)
	assert.Equal(t, "123123-K", rut.String())
	assert.False(t, rut.IsValid())

	rut, err = parseRUT("1231231-8", false)
	assert.NoError(t, err)
	assert.Equal(t, []uint8{1, 2, 3, 1, 2, 3, 1}, rut.Digits)
	assert.Equal(t, VD8, rut.VD)
	assert.Equal(t, "1231231-8", rut.String())
	assert.True(t, rut.IsValid())
}
