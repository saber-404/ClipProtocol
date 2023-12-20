package ClipProtocol

import (
	"fmt"
	"testing"
)

func TestSendStopCmd(t *testing.T) {
	var cpc ProtocolClient
	cpc.SendStopCmd(9000)
}

func TestSendData(t *testing.T) {
	var cpc ProtocolClient
	data := RandomAlphaString(20)
	fmt.Println(data)
	cpc.SendData(9000, data)
	//cpc.BetaSendData(9000, data)
}

func TestNewClient(t *testing.T) {
	NewClient().SendData(9000, "test")
	NewClient().SendStopCmd(9000)
}
