package ClipProtocol

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func Test_setFlag(t *testing.T) {
	var c ClipProtocolPacket
	fmt.Println(c.Flag)
	err := c.setFlag(0, 0, 996)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c.Flag)
}

func Test_checkPacket(t *testing.T) {
	var c ClipProtocolPacket
	p := bytes.Buffer{}
	c.ProtocolName = 0x32706330
	c.PacketID = 0
	c.PacketNum = 1
	c.PacketCount = 2
	//c.Data = []byte("2pc0")
	err := c.setFlag(1, 0, uint16(len(c.Data)))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c.Flag)
	var data = []any{c.ProtocolName, c.PacketID, c.PacketNum, c.PacketCount, c.Flag, c.Data}
	for _, v := range data {
		binary.Write(&p, binary.BigEndian, v)
	}
	pstr := hex.EncodeToString(p.Bytes())
	fmt.Println(pstr)
	packetType, ret := c.checkPacket(p.Bytes())
	fmt.Println(packetType, ret)
}

func Test_main(t *testing.T) {
	data := "2pc0"
	fmt.Println(hex.EncodeToString([]byte(data)))
}

func Test_getPacketData(t *testing.T) {
	var c ClipProtocolPacket
	p := bytes.Buffer{}
	c.ProtocolName = 0x32706330
	c.PacketID = 0
	c.PacketNum = 1
	c.PacketCount = 2
	c.Data = []byte("2pc06565")
	err := c.setFlag(0, 0, uint16(len(c.Data)))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c.Flag)
	var data = []any{c.ProtocolName, c.PacketID, c.PacketNum, c.PacketCount, c.Flag, c.Data}
	for _, v := range data {
		binary.Write(&p, binary.BigEndian, v)
	}
	fmt.Println(hex.EncodeToString(p.Bytes()))
	pdata, err := c.getPacketData(p.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(pdata))
}

func Test_genStopCmdPacket(t *testing.T) {
	var c ClipProtocolPacket
	PacketID := uint64(time.Now().Unix())
	println(PacketID)
	stopCmdPacket := c.genStopCmdPacket(PacketID)
	packetType, ret := c.checkPacket(stopCmdPacket)
	println(hex.EncodeToString(stopCmdPacket))
	fmt.Println(packetType, ret)
}

func Test_genDataPacket(t *testing.T) {
	var c ClipProtocolPacket
	PacketID := uint64(time.Now().Unix())
	println(PacketID)
	err, i := c.genDataPacket(PacketID, 0, 1, DataPacketFlag, []byte("2pc06565"))
	if err != nil {
		fmt.Println(err)
		return
	}
	packetType, ret := c.checkPacket(i)
	fmt.Println(packetType, ret)
	fmt.Println(hex.EncodeToString(i))
}
