package ClipProtocol

import (
	"fmt"
	"testing"
)

func TestClipProtocolServer_StartServer(t *testing.T) {
	server := NewServer("0.0.0.0:9000")
	server.Run(HandlerFunc)
}

func HandlerFunc(str string) {
	fmt.Println(str)
}
