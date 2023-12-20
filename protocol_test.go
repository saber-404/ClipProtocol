package ClipProtocol

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"testing"
	"time"
)

func Test_setFlag(t *testing.T) {
	var c ProtocolPacket
	fmt.Println(c.Flag)
	err := c.setFlag(0, 0, 996)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c.Flag)
}

func Test_checkPacket(t *testing.T) {
	var c ProtocolPacket
	p := bytes.Buffer{}
	c.ProtocolName = 0x32706330
	c.PacketID = 0
	c.PacketNum = 1
	c.PacketCount = 2
	c.Data = []byte("2pc0")
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
	pstr := hex.EncodeToString(p.Bytes())
	fmt.Println(pstr)
	err, PacketID, PacketNum, PacketCount, Flag, Data := c.parsePacket(p.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(PacketID, PacketNum, PacketCount, Flag, Data)
}

func Test_main(t *testing.T) {
	/*	data := "2pc0"
		fmt.Println(hex.EncodeToString([]byte(data)))*/
	//fmt.Println(10 / 3)
	//fmt.Println(10 % 2)
	//
	//a, b := 10, 3
	//result := a / b
	//ceilResult := math.Ceil(float64(result))
	//
	//fmt.Printf("a / b = %d\n", result)
	//fmt.Printf("a / b 的向上取整结果 = %.2f\n", ceilResult)

	//m := make(map[int]string, 9)
	//m[1] = "z"
	//m[5] = "b"
	//m[1] = "d"
	//fmt.Println(len(m))
	// 获取本机的IP地址
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取IP地址失败:", err)
		return
	}

	// 遍历所有IP地址
	for _, addr := range addrs {
		// 检查IP地址是否是IPv4或IPv6
		ip, ok := addr.(*net.IPNet)
		if ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				// 打印IPv4地址
				fmt.Println("IP地址:", ip.IP)
				// 获取广播地址
				broadcast := getBroadcastAddress(ip)
				fmt.Println("广播地址:", broadcast)
			}
		}
	}
}

// 获取广播地址
func getBroadcastAddress(ip *net.IPNet) net.IP {
	// 将IP地址转换为4字节表示
	ip = &net.IPNet{IP: ip.IP.To4(), Mask: ip.Mask}
	// 计算广播地址
	broadcast := make(net.IP, len(ip.IP))
	for i, b := range ip.IP {
		broadcast[i] = b | ^ip.Mask[i]
	}
	return broadcast
}

func Test_getPacketData(t *testing.T) {
	var c ProtocolPacket
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
	var c ProtocolPacket
	PacketID := uint64(time.Now().Unix())
	fmt.Println("PacketID", PacketID)
	stopCmdPacket := c.genStopCmdPacket(PacketID) // PacketID 0 1 0x8000
	err, PacketID, PacketNum, PacketCount, Flag, _ := c.parsePacket(stopCmdPacket)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(PacketID, PacketNum, PacketCount, Flag)
}

func Test_genDataPacket(t *testing.T) {
	var c ProtocolPacket
	PacketID := uint64(time.Now().Unix())
	fmt.Println("PacketID", PacketID)
	err, i := c.genDataPacket(PacketID, 0, 1, DataFlag, []byte("2pc06565"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err, PacketID, PacketNum, PacketCount, _, data := c.parsePacket(i)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(PacketID, PacketNum, PacketCount, string(data), hex.EncodeToString(data))
}

func Test_genPacketSliceFromString(t *testing.T) {
	strData := strings.Repeat("2pc0", 500)
	var c ProtocolPacket
	var packetID = uint64(time.Now().Unix())
	packets, err := c.genPacketSliceFromString(packetID, strData)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range packets {
		packetID, err := c.getPacketID(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		packetNum, packetCount, err := c.getPacketNumAndCount(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("PacketID: %d, PacketNum: %d, PacketCount: %d\n", packetID, packetNum, packetCount)
	}
}

func randString(MetaString string, length int) string {
	byteofstr := []byte(MetaString)
	var result []byte
	for i := 0; i < length; i++ {
		rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
		result = append(result, byteofstr[rand.Intn(len(byteofstr))])
	}
	return string(result)
}

func RandomAlphaString(length int) string {
	MetaString := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return randString(MetaString, length)
}

//func Test_genStringFromPacketSlice(t *testing.T) {
//	strData := RandomAlphaString(2000)
//	var c ProtocolPacket
//	var packetID = uint64(time.Now().Unix())
//	packets, _ := c.genPacketSliceFromString(packetID, strData)
//	str, err := c.genStringDataOfFromPacketSlice(packets)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if str == strData {
//		fmt.Println("OK")
//	} else {
//		fmt.Println("FAIL")
//	}
//	return
//}
