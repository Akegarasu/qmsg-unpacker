package main

import (
	"fmt"
	"github.com/Akegarasu/qmsg-unpacker/qqmsg"
	"os"
)

func main() {
	b, err := os.ReadFile("a_example_msg_content")
	if err != nil {
		return
	}
	msg := qqmsg.Unpack(b)
	fmt.Println(msg.SenderNickname)
	fmt.Println(msg.Header.Time)
	msgPrintable := qqmsg.EncodeMsg(msg)
	fmt.Println(msgPrintable)
}
