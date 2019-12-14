package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Unpack(t *testing.T) {
	unpacker := RepCharUnpacker{}

	res, err := unpacker.Unpack("")
	assert.NoError(t, err)
	assert.Equal(t, "", res)

	res, err = unpacker.Unpack("abcd")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", res)

	res, err = unpacker.Unpack("a2")
	assert.NoError(t, err)
	assert.Equal(t, "aa", res)

	res, err = unpacker.Unpack("45")
	assert.Error(t, err)
	assert.Equal(t, "", res)

	res, err = unpacker.Unpack("a0")
	assert.Error(t, err)
	assert.Equal(t, "", res)

	res, err = unpacker.Unpack("a02")
	assert.Error(t, err)
	assert.Equal(t, "", res)

	res, err = unpacker.Unpack("a4bc2d5e")
	assert.NoError(t, err)
	assert.Equal(t, "aaaabccddddde", res)

	res, err = unpacker.Unpack("ab11c")
	assert.NoError(t, err)
	assert.Equal(t, "abbbbbbbbbbbc", res)

	res, err = unpacker.Unpack(`qwe\4\5`)
	assert.NoError(t, err)
	assert.Equal(t, `qwe45`, res)

	res, err = unpacker.Unpack(`qwe\45`)
	assert.NoError(t, err)
	assert.Equal(t, `qwe44444`, res)

	res, err = unpacker.Unpack(`qwe\\5`)
	assert.NoError(t, err)
	assert.Equal(t, `qwe\\\\\`, res)

	res, err = unpacker.Unpack(`qwe\`)
	assert.Error(t, err)
	assert.Equal(t, "", res)

	res, err = unpacker.Unpack(`qwХа4\56ё`)
	assert.NoError(t, err)
	assert.Equal(t, "qwХаааа555555ё", res)
}
