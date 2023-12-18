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

type ClipProtocolServer struct {
	Conn          *net.UDPConn
	CPP           ClipProtocolPacket
	DataMap       map[uint64]*HandlerData
	quit          chan int
	dataTransOver chan int
	Result        string
}

type HandlerData struct {
	PacketSlice []Packet
	Check       map[uint64]bool
}

func NewServer(address string) *ClipProtocolServer {
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
	return &ClipProtocolServer{
		Conn:          conn,
		CPP:           ClipProtocolPacket{},
		DataMap:       map[uint64]*HandlerData{},
		quit:          make(chan int, 1),
		dataTransOver: make(chan int, 1),
	}
}

func (CPS *ClipProtocolServer) Run(handler func(string)) {
	defer CPS.Conn.Close()
	go CPS.clearMem(ClearTimeOut)
	for {
		select {
		case <-CPS.quit:
			//fmt.Println("quit")
			return

		case <-CPS.dataTransOver:
			//fmt.Println("dataTransOver")
			handler(CPS.Result)
			CPS.Result = ""
		default:
			CPS.run()
		}
	}
}

func (CPS *ClipProtocolServer) HandlerDataPacket(packet Packet) {
	packetType, ret := CPS.CPP.checkPacket(packet)
	if packetType && ret { //处理终止命令包
		CPS.quit <- 1
		return
	}
	if !packetType && ret { //处理数据包
		packetID, _ := CPS.CPP.getPacketID(packet)
		packetNum, packetCount, _ := CPS.CPP.getPacketNumAndCount(packet)
		v, exists := CPS.DataMap[packetID]
		if !exists {
			CPS.DataMap[packetID] = &HandlerData{
				PacketSlice: make([]Packet, packetCount),
				Check:       make(map[uint64]bool, packetCount),
			}
			CPS.DataMap[packetID].PacketSlice[packetNum] = packet
			CPS.DataMap[packetID].Check[packetNum] = true
		} else {
			// 判断packetNum的包是否已经接收过
			if v.Check[packetNum] {
				return
			} else {
				v.Check[packetNum] = true
				v.PacketSlice[packetNum] = packet
			}
		}
		if uint64(len(CPS.DataMap[packetID].Check)) == packetCount {
			str, err := CPS.CPP.genStringFromPacketSlice(CPS.DataMap[packetID].PacketSlice)
			if err != nil {
				fmt.Println(err)
				return
			}
			CPS.dataTransOver <- 1
			CPS.Result = str
			//fmt.Println(str)
			//fmt.Println(len(str))
		}
	}
}

func (CPS *ClipProtocolServer) run() {
	buf := make([]byte, MaxClipProtocolPacketSize)
	n, addr, err := CPS.Conn.ReadFromUDP(buf)
	fmt.Println(addr) //todo 处理addr
	if err != nil {
		fmt.Println(err)
		return
	}
	packet := buf[:n]
	CPS.HandlerDataPacket(packet)
}

// 清理 ClearTimeOut 分钟以前接收的数据
func (CPS *ClipProtocolServer) clearMem(timeout time.Duration) {
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
