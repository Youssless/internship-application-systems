package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/ipv4"

	"golang.org/x/net/icmp"
)

func main() {
	// default destination is www.google.com
	destName := flag.String("dest", "www.google.com", "string")
	numOfPings := flag.Int("n", 20, "int")
	flag.Parse()

	var i int = 0
	for i < *numOfPings {
		ping(*destName)
		i++
	}
}

// ICMP packet
func packet(ipversion icmp.Type) *icmp.Message {
	packt := icmp.Message{
		Type: ipversion,
		Code: 0, // 0 = echo
		// payload
		Body: &icmp.Echo{
			ID:   os.Getpid(),
			Seq:  0,
			Data: []byte(""),
		},
	}
	return &packt
}

func ping(destName string) {
	fmt.Println("Destination Name", destName)
	// icmp listener connection
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println(err)
		return
	}
	listenerAddr := conn.LocalAddr()

	// encode the packet with the message into bytes
	message := packet(ipv4.ICMPTypeEcho)
	packetBytes, err := message.Marshal(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// get the destination address
	destinationAddr, err := net.ResolveIPAddr("ip4:icmp", destName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// logging info to the terminal
	messageParsed, _ := icmp.ParseMessage(1, packetBytes)
	fmt.Printf("[From: %s To: %s] => Message: %v", listenerAddr, destinationAddr, *messageParsed)

	startTime := time.Now()

	// send to destination address
	_, err = conn.WriteTo(packetBytes, destinationAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// get the reply from the listener
	reply := make([]byte, 1500)
	conn.ReadFrom(reply)

	endTime := time.Since(startTime)

	// logging info to the terminal
	replyMessage, _ := icmp.ParseMessage(1, reply)
	fmt.Printf(" [Time: %s] => Reply: %v\n", endTime, *replyMessage)
	conn.Close()
}
