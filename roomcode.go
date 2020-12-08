package goamongus

import (
	"errors"
	"strings"
)

// Without static class functions this is a little messy.
// So expect the API to break when I decide on a cleaner interface.

var CharSet = []int32{
	65: 0x19, 0x15, 0x13, 0x0a, 0x08,
	0x0b, 0x0c, 0x0d, 0x16, 0x0f,
	0x10, 0x06, 0x18, 0x17, 0x12,
	0x07, 0x00, 0x03, 0x09, 0x04,
	0x0e, 0x14, 0x01, 0x02, 0x05,
	0x11}

func Encode(code string) (int32, error) {
	code = strings.ToUpper(code)

	if len(code) == 6 {
		return encodeV2(code), nil
	}

	return 0, errors.New("Invalid room code length, expected 6 characters: " + code)
}

func encodeV2(codeStr string) int32 {
	b1 := CharSet[codeStr[0]]
	b2 := CharSet[codeStr[1]]
	b3 := CharSet[codeStr[2]]
	b4 := CharSet[codeStr[3]]
	b5 := CharSet[codeStr[4]]
	b6 := CharSet[codeStr[5]]

	lsb := (b1 + 26*b2) & 0x3FF
	msb := (b3 + 26*(b4+26*(b5+26*b6)))

	return lsb | ((msb << 10) & 0x3FFFFC00) | -0x80000000
}
