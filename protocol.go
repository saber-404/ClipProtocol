package ClipProtocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	MaxClipProtocolPacketDataSize        = 994
	MaxClipProtocolPacketSize            = 1024
	MinClipProtocolPacketSize            = 30
	PacketProtocolName            uint32 = 0x32706330
	PacketProtocolNameString             = "2pc0"
	StopCmdPacketFlag             uint16 = 0x8000
	DataFlag                      uint16 = 0
	StopCmdFlag                   uint16 = 1
)

type Packet []byte
type DataOfPacket []byte

type ProtocolPacket struct {
	ProtocolName uint32 // 协议名 4字节   "2pc0" 0x32706330 846226224
	PacketID     uint64 // 发送数据时间戳 8字节
	PacketNum    uint64 // 第几个包  8字节	最小 0
	PacketCount  uint64 // 总包数	8字节 	最小 1
	Flag         uint16 // 传输类型 1位 + 占位 5 位 + 数据长度 10位(最大数据量1024字节) 2字节
	Data         []byte // 实际数据最大 994字节
}

/*
配置Flag数据到 ProtocolPacket.Flag
trantype 传输类型 0 -> 数据包 1 -> 终止命令(此时dataLen为0 Flag固定为32768)
other 占位
dataLen 数据长度 最大994
*/
func (cpp *ProtocolPacket) setFlag(trantype, other, dataLen uint16) (err error) {
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

//解析数据包
func (cpp *ProtocolPacket) parsePacket(packet []byte) (err error, PacketID, PacketNum, PacketCount uint64, Flag uint16, Data []byte) {
	PacketLength := uint16(len(packet))
	if PacketLength < MinClipProtocolPacketSize ||
		PacketLength > MaxClipProtocolPacketSize {
		err = errors.New("PacketLength is too wrong")
		return
	}
	s := string(packet[0:4])
	if s != PacketProtocolNameString {
		err = errors.New("PacketProtocolName is not match")
	}
	binary.Read(bytes.NewBuffer(packet[4:12]), binary.BigEndian, &PacketID)
	binary.Read(bytes.NewBuffer(packet[12:20]), binary.BigEndian, &PacketNum)
	binary.Read(bytes.NewBuffer(packet[20:28]), binary.BigEndian, &PacketCount)
	binary.Read(bytes.NewBuffer(packet[28:MinClipProtocolPacketSize]), binary.BigEndian, &Flag)
	if PacketCount <= PacketNum {
		err = errors.New("PacketCount <= PacketNum")
		return
	}
	if Flag == StopCmdPacketFlag {
		if PacketLength != MinClipProtocolPacketSize {
			err = errors.New("StopCmdPacket Length is not match")
			return
		} else {
			err = nil
			return
		}
	} else if Flag > StopCmdPacketFlag {
		err = errors.New("wrong Packet Flag")
		return
	} else {
		if Flag != PacketLength-MinClipProtocolPacketSize {
			err = errors.New("wrong Packet Flag")
			return
		}
	}
	Data = packet[MinClipProtocolPacketSize:]
	err = nil
	return
}

//获取数据包中的数据
func (cpp *ProtocolPacket) getPacketData(packet []byte) (data []byte, err error) {
	err, _, _, _, _, data = cpp.parsePacket(packet)
	return
}

//获取数据包ID
func (cpp *ProtocolPacket) getPacketID(packet []byte) (packetID uint64, err error) {
	err, packetID, _, _, _, _ = cpp.parsePacket(packet)
	return
}

//获取数据包 次序 总数
func (cpp *ProtocolPacket) getPacketNumAndCount(packet []byte) (packetNum, packetCount uint64, err error) {
	err, _, packetNum, packetCount, _, _ = cpp.parsePacket(packet)
	return
}

//生成结束命令的命令包
//packetID 命令包ID 可以使用时间戳
func (cpp *ProtocolPacket) genStopCmdPacket(packetID uint64) []byte {
	err, i := cpp.genDataPacket(packetID, 0, 1, StopCmdFlag, nil)
	if err != nil {
		return nil
	}
	return i
}

//生成数据包
//packetID 命令包ID 可以使用时间戳
//PacketNum 第几个包
//PacketCount 总包数
//trantype 数据包类型 DataFlag -> 数据包 DataStopCmdFlag -> 结束命令
//data 包内数据最大 MaxClipProtocolPacketDataSize
func (cpp *ProtocolPacket) genDataPacket(packetID, PacketNum, PacketCount uint64, trantype uint16, data []byte) (err error, packet []byte) {
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
	var packetdata = []any{PacketProtocolName, packetID, PacketNum, PacketCount, cpp.Flag, data}
	for _, v := range packetdata {
		binary.Write(&p, binary.BigEndian, v)
	}
	err = nil
	packet = p.Bytes()
	return
}

//从字符串生成数据包切片[]Packet
func (cpp *ProtocolPacket) genPacketSliceFromString(packetID uint64, strData string) (packetSlice []Packet, err error) {
	//	字符串转[]byte
	data := []byte(strData)
	l := uint64(len(data))
	//	计算出PacketNum PacketCount
	PacketCount := l / MaxClipProtocolPacketDataSize
	if l%MaxClipProtocolPacketDataSize != 0 {
		PacketCount += 1
	}
	var (
		i uint64
	)
	for i = 0; i < l; i += MaxClipProtocolPacketDataSize {
		var p Packet
		if i+MaxClipProtocolPacketDataSize < l {
			err, p = cpp.genDataPacket(packetID, i/MaxClipProtocolPacketDataSize, PacketCount, DataFlag, data[i:i+MaxClipProtocolPacketDataSize])
			if err != nil {
				return nil, err
			}
		} else {
			err, p = cpp.genDataPacket(packetID, i/MaxClipProtocolPacketDataSize, PacketCount, DataFlag, data[i:])
			if err != nil {
				return nil, err
			}
		}
		packetSlice = append(packetSlice, p)
	}
	return
}

//从[]Packet生成字符串
func (cpp *ProtocolPacket) genStringDataOfFromPacketSlice(packetSlice []DataOfPacket) (strData string) {
	var (
		byteOfStr []byte
	)
	for _, v := range packetSlice {
		byteOfStr = append(byteOfStr, v...)
	}
	return string(byteOfStr)
	//var (
	//	byteOfStr     []byte
	//	FirstPacketID uint64
	//)
	//FirstPacketID, err = cpp.getPacketID(packetSlice[0])
	//l := uint64(len(packetSlice))
	//for i, packet := range packetSlice {
	//	var (
	//		PacketID    uint64
	//		PacketData  []byte
	//		PacketNum   uint64
	//		PacketCount uint64
	//	)
	//	PacketID, err = cpp.getPacketID(packet)
	//	if err != nil {
	//		return "", err
	//	}
	//	if FirstPacketID != PacketID {
	//		return "", errors.New("the packetSlice has different PacketID")
	//	}
	//	PacketNum, PacketCount, err = cpp.getPacketNumAndCount(packet)
	//	if err != nil {
	//		return "", err
	//	}
	//	if uint64(i) != PacketNum {
	//		return "", errors.New("the packetSlice order is wrong")
	//	}
	//	if l != PacketCount {
	//		return "", errors.New("the packetSlice length is wrong")
	//	}
	//	PacketData, err = cpp.getPacketData(packet)
	//	if err != nil {
	//		return "", err
	//	}
	//	byteOfStr = append(byteOfStr, PacketData...)
	//}
	//return string(byteOfStr), nil
}

func (cpp *ProtocolPacket) genReplyClientMessagePacket(packetID uint64) (err error, p Packet) {
	err, p = cpp.genDataPacket(packetID, 0, 1, DataFlag, []byte{1})
	if err != nil {
		return err, nil
	}
	return
}
