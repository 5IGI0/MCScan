package main

import (
	"encoding/binary"
	"errors"
)

type MCString string

func (s *MCString) SizeOf() int {
	l := VarUInt(len(*s))
	return l.SizeOf() + len(*s)
}

func (s *MCString) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	var l VarUInt
	buf, err := l.Unpack(buf, order)
	if err != nil {
		return []byte{}, err
	}
	if len(buf) < int(l) {
		return []byte{}, errors.New("buffer too small")
	}
	*s = MCString(buf[:l])
	return buf[l:], nil
}

func (s *MCString) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	l := VarUInt(len(*s))
	buf, err := l.Pack(buf, order)
	if err != nil {
		return []byte{}, err
	}
	return buf[copy(buf, *s):], nil
}
