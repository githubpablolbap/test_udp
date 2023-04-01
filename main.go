package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, _ *http.Request) {

	fmt.Fprintf(w, "Hello, there\n")
}

func ReadUdp() {
	udpServer, err := net.ListenPacket("udp", ":55555")
	if err != nil {
		log.Fatal(err)
	}
	//
	defer udpServer.Close()
	for {
		buf := make([]byte, 1024)
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		udpServer.WriteTo(buf[:n], addr)
		// go response(udpServer, addr, buf)
	}
}

// func response(udpServer net.PacketConn, addr net.Addr, buf []byte) {
// 	time := time.Now().Format(time.ANSIC)
// 	responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf))
// 	udpServer.WriteTo([]byte(responseStr), addr)
// 	// c.WriteTo(packet[:n], addr)
// }

func main() {
	//
	// go ReadUdp()
	//
	http.HandleFunc("/", HelloHandler)
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":55555", nil))
}
