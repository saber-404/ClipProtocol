package ClipProtocol

import (
	"fmt"
	"log"
	"net"
	"time"
)

type ProtocolClient struct {
	CPP      ProtocolPacket
	PacketID uint64
	conn     net.Conn
}

func NewClient() *ProtocolClient {
	return &ProtocolClient{}
}

func (cpc *ProtocolClient) SendData(port int, data string) {
	address := fmt.Sprintf("255.255.255.255:%d", port)
	conn, err := net.Dial("udp", address)
	if err != nil {
		log.Fatal(err)
		return
	}
	//cpc.conn = conn
	defer cpc.conn.Close()
	cpc.PacketID = uint64(time.Now().Unix())
	packets, err := cpc.CPP.genPacketSliceFromString(cpc.PacketID, data)
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
	//c := make(chan int, 1)
	//for {
	//	select {
	//	case c <- 1:
	//		return
	//	default:
	//		cpc.GetReplyData(c)
	//	}
	//}
}

//func (cpc *ProtocolClient) BetaSendData(port int, data string) {
//	address := fmt.Sprintf("255.255.255.255:%d", port)
//	conn, err := net.Dial("udp", address)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	cpc.conn = conn
//	defer cpc.conn.Close()
//	c := make(chan int, 1)
//	go cpc.GetReplyData(c)
//	cpc.PacketID = uint64(time.Now().Unix())
//	packets, err := cpc.CPP.genPacketSliceFromString(cpc.PacketID, data)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	for _, packet := range packets {
//		n, err := conn.Write(packet)
//		if err != nil {
//			return
//		}
//		fmt.Println("Sent:", n, "bytes")
//	}
//}

func (cpc *ProtocolClient) SendStopCmd(port int) {
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

//func (cpc *ProtocolClient) GetReplyData(c chan int) {
//	for {
//		buf := make([]byte, 1024)
//		n, err := cpc.conn.Read(buf)
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//		fmt.Println("GetReplyData:", n, "bytes")
//		packetID, err := cpc.CPP.getPacketID(buf[:n])
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//		if packetID == cpc.PacketID {
//			c <- 1
//			fmt.Println("GetReplyData:", packetID)
//			break
//		}
//	}
//}
