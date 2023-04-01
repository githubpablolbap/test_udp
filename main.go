package main

//go build -ldflags="-s -w" server.go
import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/", HandleEntry)
	fmt.Println("...server started...")
	// http.ListenAndServe(":"+GoPort("22222"), nil)
	http.ListenAndServe(":22222", nil)
}

func HandleEntry(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/ws/" {
		upgrader := websocket.Upgrader{ReadBufferSize: 128, WriteBufferSize: 1024}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err == nil {
			// var client *Client = &Client{conn, 0, 0, 100, [20]byte{0}}
			// // fmt.Println(client)
			go HandleClient(conn)
			// *VALIDATE = append(*VALIDATE, client)
			// validate += 1
			// return
		} else {
			conn.Close()
			return
		}
	} else {
		fmt.Fprintf(w, "Hello, there\n")
		// fmt.Println(w, "/")
		return
	}
}
func GoPort(p string) string {
	pt := os.Getenv("PORT")
	if pt == "" {
		pt = p
	}
	return pt
}

func HandleClient(conn *websocket.Conn) {
	//
	var msg []byte
	var err error
	// var conn *websocket.Conn = client.conn
	// var nick [20]byte = [20]byte{}
	//
	for {
		_, msg, err = conn.ReadMessage()
		if err == nil {
			err = conn.WriteMessage(2, msg[:])
			if err != nil {
				fmt.Println("error SendData():", err)
			}
		} else {
			conn.Close()
			return
		}
	}
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net"
// 	"net/http"
// )

// func HelloHandler(w http.ResponseWriter, _ *http.Request) {

// 	fmt.Fprintf(w, "Hello, there\n")
// }

// func ReadUdp() {
// 	udpServer, err := net.ListenPacket("udp", ":55555")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//
// 	defer udpServer.Close()
// 	for {
// 		buf := make([]byte, 1024)
// 		n, addr, err := udpServer.ReadFrom(buf)
// 		if err != nil {
// 			continue
// 		}
// 		udpServer.WriteTo(buf[:n], addr)
// 		// go response(udpServer, addr, buf)
// 	}
// }

// // func response(udpServer net.PacketConn, addr net.Addr, buf []byte) {
// // 	time := time.Now().Format(time.ANSIC)
// // 	responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf))
// // 	udpServer.WriteTo([]byte(responseStr), addr)
// // 	// c.WriteTo(packet[:n], addr)
// // }

// func main() {
// 	//
// 	go ReadUdp()
// 	//
// 	http.HandleFunc("/", HelloHandler)
// 	log.Println("Listening...")
// 	log.Fatal(http.ListenAndServe(":55555", nil))
// }
