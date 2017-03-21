package main

import (
	"github.com/clairezz/golab/ccp/ccp"
)

// 遵循 Channel Closing Principle, 并优雅地关闭channel
func main() {
//	ccp.Sender2Receivers()
//	ccp.Senders2Receiver()
//	ccp.Senders2Receivers()
	ccp.Senders2Receivers2()
}
