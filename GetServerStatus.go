package main

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-restruct/restruct"
)

type ServerQueryResult struct {
	Address        string
	AddrAliases    []string
	StatusRaw      string
	StatusObj      MCServerStatus
	NormalizedDesc string
	Mods           string
	Favicon        []byte
}

func GetServerStatus(serv_addr string) (ServerQueryResult, error) {
	var err error
	var ret ServerQueryResult
	var connaddr string

	ret.Address = serv_addr
	if connaddr, ret.AddrAliases, err = ResolvMcAddr(serv_addr); err != nil {
		return ret, err
	}

	status_payload, err := requestStatus(serv_addr, connaddr)
	err = AnalyzeServerStatus(status_payload, &ret)

	return ret, err
}

// what a mess
func requestStatus(srvaddr string, connaddr string) (string, error) {
	var err error
	var ret string
	dialer := net.Dialer{Timeout: time.Second * SCAN_CONN_TIMEOUT}
	conn, err := dialer.Dial("tcp", connaddr)
	if err != nil {
		return ret, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second * SCAN_READ_TIMEOUT))

	sconnaddr := strings.Split(connaddr, ":")
	p, _ := strconv.Atoi(sconnaddr[len(sconnaddr)-1])
	a := McHandShake{
		ProtoVer:  MC_STATUS_VER,
		ServAddr:  MCString(srvaddr),
		ServPort:  uint16(p),
		NextState: 1}
	b := McPacketBase{}
	b.Data, _ = restruct.Pack(binary.BigEndian, &a)
	b.Length = VarUInt(len(b.Data)) + VarUInt(b.PacketId.SizeOf())
	buff, _ := restruct.Pack(binary.BigEndian, &b)
	if _, err := conn.Write(buff); err != nil {
		return ret, err
	}

	b.Data = nil
	b.Length = 1
	b.PacketId = 0
	buff, _ = restruct.Pack(binary.BigEndian, &b)
	if _, err := conn.Write(buff); err != nil {
		return ret, err
	}

	var aa [9]byte
	n, err := io.ReadFull(conn, aa[:])

	if err != nil {
		return ret, err
	}
	if n != 9 {
		return ret, errors.New("connection closed")
	}

	var l VarUInt
	if buff, err = l.Unpack(aa[:], binary.BigEndian); err != nil {
		return ret, err
	}
	if l > 100000 {
		return ret, errors.New("status too big")
	}
	if l < VarUInt(len(buff)) {
		return ret, errors.New("???")
	}
	buff = make([]byte, 9+int(l)-len(buff))
	copy(buff, aa[:])
	if _, err := io.ReadFull(conn, buff[9:]); err != nil {
		return ret, err
	}

	var r VarUInt
	var status_payload MCString
	buff, _ = r.Unpack(buff, binary.BigEndian)
	buff, _ = r.Unpack(buff, binary.BigEndian)
	_ = r
	buff, err = status_payload.Unpack(buff, binary.BigEndian)

	return string(status_payload), err
}
