package main

//go build -ldflags="-s -w" server.go
import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/", HandleEntry)
	fmt.Println("...server started...")
	// http.ListenAndServe(":"+GoPort("22222"), nil)
	http.ListenAndServe(":8080", nil)
}

func HandleEntry(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/ws/" {
		// fmt.Fprintf(w, "Hello, there .../ws/\n")
		// upgrader := websocket.Upgrader{ReadBufferSize: 128, WriteBufferSize: 1024}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err == nil {
			go HandleClient(conn)
			// return
		} else {
			conn.Close()
			// return
		}
	} else {
		fmt.Fprintf(w, "Hello, there .../\n")
		// fmt.Fprintf(w,"%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
		// return
	}
}

// func GoPort(p string) string {
// 	pt := os.Getenv("PORT")
// 	if pt == "" {
// 		pt = p
// 	}
// 	return pt
// }

func HandleClient(conn *websocket.Conn) {
	//
	// var msg []byte
	var err error
	// fmt.Println(msg[:1])
	// var conn *websocket.Conn = client.conn
	// var nick [20]byte = [20]byte{}
	//
	for {
		// _, msg, err = conn.ReadMessage()
		_, _, err = conn.ReadMessage()
		if err == nil {
			err = conn.WriteMessage(2, []byte{0, 0, 0, 0, 0})
			if err != nil {
				fmt.Println("error SendData():", err)
			}
		} else {
			conn.Close()
			break
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
