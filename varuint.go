package main

import (
	"encoding/binary"
	"errors"
)

type VarUInt uint64

func (s *VarUInt) SizeOf() int {
	var i int
	for i = 1; (*s >> (7 * i)) != 0; i++ {
	}
	return i
}

func (s *VarUInt) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	*s = 0
	for i, l := 0, len(buf); i < l && i < 8; i++ {
		*s |= VarUInt(VarUInt(buf[i]&0x7F) << (7 * i))

		if (buf[i] & 0x80) == 0 {
			return buf[i+1:], nil
		}
	}
	return []byte{}, errors.New("unterminated varint")
}

func (s *VarUInt) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	var i int
	buf[0] = byte(*s) | 0x80

	for i = 1; (*s >> (7 * i)) != 0; i++ {
		buf[i] |= byte(*s>>(7*i)) | 0x80
	}
	buf[i-1] &= 0x7F
	return buf[i:], nil
}
