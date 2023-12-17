package ClipProtocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	MaxClipProtocolPacketDataSize               = 994
	MaxClipProtocolPacketSize                   = 1024
	MinClipProtocolPacketSize                   = 30
	ClipProtocolPacketProtocolName       uint32 = 0x32706330
	ClipProtocolPacketProtocolNameString        = "2pc0"
	ClipProtocolPacketStopCmdFlag        uint16 = 0x8000
	DataPacketFlag                       uint16 = 0
	StopCmdFlag                          uint16 = 1
)

type ClipProtocolPacket struct {
	ProtocolName uint32 // 协议名 4字节   "2pc0" 0x32706330 846226224
	PacketID     uint64 // 发送数据时间戳 8字节
	PacketNum    uint64 // 第几个包  8字节	最小 0
	PacketCount  uint64 // 总包数	8字节 	最小 1
	Flag         uint16 // 传输类型 1位 + 占位 5 位 + 数据长度 10位(最大数据量1024字节) 2字节
	Data         []byte // 实际数据最大 994字节
}

/*
配置Flag数据到 ClipProtocolPacket.Flag
trantype 传输类型 0 -> 数据包 1 -> 终止命令(此时dataLen为0 Flag固定为32768)
other 占位
dataLen 数据长度 最大994
*/
func (cpp *ClipProtocolPacket) setFlag(trantype, other, dataLen uint16) (err error) {
	if other != 0 {
		err = fmt.Errorf("other is not 0")
		return
	}
	if trantype == 1 { // 终止命令
		cpp.Flag = trantype << 15 //32768
		return
	} else if trantype == 0 { //传输数据
		if dataLen > MaxClipProtocolPacketDataSize {
			err = fmt.Errorf("dataLen is too long, maxis 994 Bytes")
			return
		}
		cpp.Flag = dataLen
		return err
	} else {
		err = fmt.Errorf("trantype is not 0 or 1")
		return
	}
}

//检测数据包是否合法
//ret 是否合法
//packetType True -> 类型命令包 False -> 数据包（仅ret = True时)或不合法的数据包
func (cpp *ClipProtocolPacket) checkPacket(packet []byte) (packetType, ret bool) {
	PacketLength := uint16(len(packet))
	if PacketLength < MinClipProtocolPacketSize {
		ret = false
		return
	}
	s := string(packet[0:4])
	if s != ClipProtocolPacketProtocolNameString {
		ret = false
	}
	var (
		PacketNum   uint64
		PacketCount uint64
		Flag        uint16
	)
	binary.Read(bytes.NewBuffer(packet[12:20]), binary.BigEndian, &PacketNum)
	binary.Read(bytes.NewBuffer(packet[20:28]), binary.BigEndian, &PacketCount)
	binary.Read(bytes.NewBuffer(packet[28:MinClipProtocolPacketSize]), binary.BigEndian, &Flag)
	if PacketCount <= PacketNum {
		ret = false
		return
	}
	if Flag == ClipProtocolPacketStopCmdFlag {
		return true, PacketLength == MinClipProtocolPacketSize
	} else if Flag > ClipProtocolPacketStopCmdFlag {
		ret = false
		return
	} else {
		ret = Flag == PacketLength-MinClipProtocolPacketSize
		return
	}
}

//获取数据包中的数据
func (cpp *ClipProtocolPacket) getPacketData(packet []byte) (data []byte, err error) {
	packetType, ret := cpp.checkPacket(packet)
	if !ret {
		return nil, errors.New("packet is error")
	} else {
		if packetType { // 命令包
			return nil, nil
		} else { //数据包
			return packet[MinClipProtocolPacketSize:], nil
		}
	}
}

//获取数据包ID
func (cpp *ClipProtocolPacket) getPacketID(packet []byte) (packetID uint64, err error) {
	_, ret := cpp.checkPacket(packet)
	if !ret {
		return 0, errors.New("packet is error")
	}
	binary.Read(bytes.NewBuffer(packet[4:12]), binary.BigEndian, &packetID)
	return
}

//获取数据包 次序 总数
func (cpp *ClipProtocolPacket) getPacketCount(packet []byte) (packetNum, packetCount uint64, err error) {
	_, ret := cpp.checkPacket(packet)
	if !ret {
		return 0, 0, errors.New("packet is error")
	}
	binary.Read(bytes.NewBuffer(packet[12:20]), binary.BigEndian, &packetNum)
	binary.Read(bytes.NewBuffer(packet[20:28]), binary.BigEndian, &packetCount)
	return
}

//生成结束命令的命令包
//packetID 命令包ID 可以使用时间戳
func (cpp *ClipProtocolPacket) genStopCmdPacket(packetID uint64) []byte {
	err, i := cpp.genDataPacket(packetID, 0, 1, StopCmdFlag, nil)
	if err != nil {
		return nil
	}
	return i
	/*	c := ClipProtocolPacket{
			ProtocolName: ClipProtocolPacketProtocolName,
			PacketID:     packetID,
			PacketNum:    0,
			PacketCount:  0,
			Flag:         ClipProtocolPacketStopCmdFlag,
		}
		p := bytes.Buffer{}
		var data = []any{c.ProtocolName, c.PacketID, c.PacketNum, c.PacketCount, c.Flag, c.Data}
		for _, v := range data {
			binary.Write(&p, binary.BigEndian, v)
		}
		return p.Bytes()*/
}

//生成数据包
//packetID 命令包ID 可以使用时间戳
//PacketNum 第几个包
//PacketCount 总包数
//trantype 数据包类型 DataPacketFlag -> 数据包 DataStopCmdFlag -> 结束命令
//data 包内数据最大 MaxClipProtocolPacketDataSize
func (cpp *ClipProtocolPacket) genDataPacket(packetID, PacketNum, PacketCount uint64, trantype uint16, data []byte) (err error, packet []byte) {
	if PacketCount < PacketNum {
		err = errors.New("the PacketCount < PacketNum, please check it")
		return
	}
	if len(data) > MaxClipProtocolPacketDataSize {
		err = errors.New("the data size > 994, please check it")
		return
	}
	err = cpp.setFlag(trantype, 0, uint16(len(data)))
	if err != nil {
		return
	}
	p := bytes.Buffer{}
	var packetdata = []any{ClipProtocolPacketProtocolName, packetID, PacketNum, PacketCount, cpp.Flag, data}
	for _, v := range packetdata {
		binary.Write(&p, binary.BigEndian, v)
	}
	err = nil
	packet = p.Bytes()
	return
}
