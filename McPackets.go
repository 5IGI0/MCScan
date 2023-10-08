package main

type McPacketBase struct {
	Length   VarUInt
	PacketId VarUInt
	Data     []byte
}

type McHandShake struct {
	ProtoVer  VarUInt
	ServAddr  MCString
	ServPort  uint16
	NextState VarUInt
}
