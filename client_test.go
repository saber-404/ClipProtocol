package ClipProtocol

import (
	"testing"
)

func TestSendStopCmd(t *testing.T) {
	var cpc ClipProtocolClient
	cpc.SendStopCmd(9000)
}

func TestSendData(t *testing.T) {
	var cpc ClipProtocolClient
	data := RandomAlphaString(20)
	cpc.SendData(9000, data)
}

func TestNewClient(t *testing.T) {
	NewClient().SendData(9000, "test")
	NewClient().SendStopCmd(9000)
}
