package ClipProtocol

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	ClearTimeOut = 3 * time.Minute
)

type ProtocolServer struct {
	Conn *net.UDPConn
	//ClientAddr    *net.UDPAddr //todo Beta
	CPP           ProtocolPacket
	DataMap       map[uint64]*HandlerData
	quit          chan int
	dataTransOver chan int
	Result        string
}

type HandlerData struct {
	PacketSlice []DataOfPacket
	Check       map[uint64]bool
}

func NewServer(address string) *ProtocolServer {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &ProtocolServer{
		Conn:          conn,
		CPP:           ProtocolPacket{},
		DataMap:       map[uint64]*HandlerData{},
		quit:          make(chan int, 1),
		dataTransOver: make(chan int, 1),
	}
}

func (CPS *ProtocolServer) Run(handler func(string)) {
	defer CPS.Conn.Close()
	go CPS.clearMem(ClearTimeOut)
	for {
		select {
		case <-CPS.quit:
			return

		case <-CPS.dataTransOver:
			handler(CPS.Result)
			CPS.Result = ""
		default:
			CPS.run()
		}
	}
}

func (CPS *ProtocolServer) HandlerDataPacket(packet Packet) {
	err, packetID, packetNum, packetCount, Flag, Data := CPS.CPP.parsePacket(packet)
	if err != nil {
		return
	}
	if Flag == StopCmdPacketFlag { //处理终止命令包
		//CPS.replyClient(packetID) //todo Beta
		CPS.quit <- 1
		return
	}
	//处理数据包
	v, exists := CPS.DataMap[packetID]
	if !exists {
		CPS.DataMap[packetID] = &HandlerData{
			PacketSlice: make([]DataOfPacket, packetCount),
			Check:       make(map[uint64]bool, packetCount),
		}
		CPS.DataMap[packetID].PacketSlice[packetNum] = Data
		CPS.DataMap[packetID].Check[packetNum] = true
	} else {
		// 判断packetNum的包是否已经接收过
		if v.Check[packetNum] {
			return
		} else {
			v.Check[packetNum] = true
			v.PacketSlice[packetNum] = Data
		}
	}
	if uint64(len(CPS.DataMap[packetID].Check)) == packetCount {
		str := CPS.CPP.genStringDataOfFromPacketSlice(CPS.DataMap[packetID].PacketSlice)
		CPS.dataTransOver <- 1
		CPS.Result = str
		delete(CPS.DataMap, packetID)
		//CPS.replyClient(packetID) //todo Beta
	}

}

func (CPS *ProtocolServer) run() {
	buf := make([]byte, MaxClipProtocolPacketSize)
	n, addr, err := CPS.Conn.ReadFromUDP(buf)
	fmt.Println(addr) //todo Beta
	//CPS.ClientAddr = addr //todo Beta
	if err != nil {
		fmt.Println(err)
		return
	}
	packet := buf[:n]
	CPS.HandlerDataPacket(packet)
}

// 清理 ClearTimeOut 分钟以前接收的数据
func (CPS *ProtocolServer) clearMem(timeout time.Duration) {
	for {
		select {
		case <-time.After(timeout):
			if len(CPS.DataMap) == 0 {
				continue
			}
			fmt.Println("清理前", len(CPS.DataMap))
			OldPacketID := uint64(time.Now().Add(timeout).Unix())
			for packetID := range CPS.DataMap {
				if packetID < OldPacketID {
					delete(CPS.DataMap, packetID)
				}
			}
			fmt.Println("清理后", len(CPS.DataMap))
		}
	}
}

//回复客户端消息
//todo Beta
//func (CPS *ProtocolServer) replyClient(packetID uint64) {
//	err, p := CPS.CPP.genReplyClientMessagePacket(packetID)
//	if err != nil {
//		return
//	}
//	fmt.Println(CPS.ClientAddr.String())
//	n, err := CPS.Conn.WriteToUDP(p, CPS.ClientAddr)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("回复客户端消息", n)
//}
