package ClipProtocol

import (
	"fmt"
	"log"
	"net"
	"time"
)

type ClipProtocolClient struct {
	CPP ClipProtocolPacket
}

func (cpc *ClipProtocolClient) SendData(port int, data string) {
	address := fmt.Sprintf("255.255.255.255:%d", port)
	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	var packetID = uint64(time.Now().Unix())
	packets, err := cpc.CPP.genPacketSliceFromString(packetID, data)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, packet := range packets {
		n, err := conn.Write(packet)
		if err != nil {
			return
		}
		fmt.Println("Sent:", n, "bytes")
	}
}

func (cpc *ClipProtocolClient) SendStopCmd(port int) {
	address := fmt.Sprintf("255.255.255.255:%d", port)
	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	packetID := time.Now().Unix()
	packet := cpc.CPP.genStopCmdPacket(uint64(packetID))
	n, err := conn.Write(packet)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("SendStopCmd:", n, "bytes")
}
