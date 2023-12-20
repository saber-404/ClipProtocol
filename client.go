package ClipProtocol

import (
	"fmt"
	"log"
	"net"
	"time"
)

type ProtocolClient struct {
	CPP ProtocolPacket
	//PacketID      uint64   //todo Beta
	//conn          net.Conn //todo Beta
	//dataTransOver chan int //todo Beta
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
	defer conn.Close()
	packetID := uint64(time.Now().Unix())
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

//Beta版本，复用链接，可双向通信
//todo Beta
//func BetaNewClient(address string) *ProtocolClient {
//	conn, err := net.Dial("udp", address)
//	if err != nil {
//		log.Fatal(err)
//		return nil
//	}
//	//fmt.Println("remote", conn.RemoteAddr())
//	//fmt.Println("local", conn.LocalAddr())
//	return &ProtocolClient{
//		conn:          conn,
//		CPP:           ProtocolPacket{},
//		dataTransOver: make(chan int, 1),
//	}
//}
//todo Beta
//func (cpc *ProtocolClient) BetaRun() {
//	//defer cpc.conn.Close()
//	for {
//		select {
//		case <-cpc.dataTransOver:
//			fmt.Println("dataTransOver")
//			return
//		default:
//			cpc.BetaGetReplyData()
//		}
//	}
//}
//todo Beta
//func (cpc *ProtocolClient) BetaSendData(data string) {
//	cpc.PacketID = uint64(time.Now().Unix())
//	packets, err := cpc.CPP.genPacketSliceFromString(cpc.PacketID, data)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	for _, packet := range packets {
//		n, err := cpc.conn.Write(packet)
//		if err != nil {
//			return
//		}
//		fmt.Println("Sent:", n, "bytes")
//	}
//	cpc.BetaRun()
//}
//todo Beta
//func (cpc *ProtocolClient) BetaSendStopCmd() {
//	packetID := time.Now().Unix()
//	packet := cpc.CPP.genStopCmdPacket(uint64(packetID))
//	n, err := cpc.conn.Write(packet)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	fmt.Println("SendStopCmd:", n, "bytes")
//	cpc.BetaRun()
//}
//todo Beta
//func (cpc *ProtocolClient) BetaGetReplyData() {
//	buf := make([]byte, 1024)
//	n, err := cpc.conn.Read(buf)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("BetaGetReplyData:", n, "bytes")
//	packetID, err := cpc.CPP.getPacketID(buf[:n])
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	if packetID == cpc.PacketID {
//		cpc.dataTransOver <- 1
//		fmt.Println("BetaGetReplyData:", packetID)
//		return
//	}
//}
