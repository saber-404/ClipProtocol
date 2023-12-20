package ClipProtocol

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestSendStopCmd(t *testing.T) {
	NewClient().SendStopCmd(9000)
}

func TestSendData(t *testing.T) {

	data := RandomAlphaString(20)
	fmt.Println(data)
	NewClient().SendData(9000, data)
}

func TestRandomOrderPacket(t *testing.T) {
	var cpc ProtocolClient
	PacketID := uint64(time.Now().Unix())
	err, packet1 := cpc.CPP.genDataPacket(PacketID, 0, 2, DataFlag, []byte("第1个数据包"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err, packet2 := cpc.CPP.genDataPacket(PacketID, 1, 2, DataFlag, []byte("第2个数据包"))
	if err != nil {
		fmt.Println(err)
		return
	}
	packets := [][]byte{packet2, packet1}
	conn, err := net.Dial("udp", "255.255.255.255:9000")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	for _, packet := range packets {
		n, err := conn.Write(packet)
		if err != nil {
			return
		}
		fmt.Println("Sent:", n, "bytes")
	}
}

//todo Beta
//func TestBeta(t *testing.T) {
//	client := BetaNewClient("192.168.0.106:9000")
//	client.BetaSendData("test")
//	client.BetaSendStopCmd()
//}
