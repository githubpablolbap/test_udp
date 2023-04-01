package main

//go build -ldflags="-s -w" server.go
import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MAX_GAMES uint32 = uint32(3)
)

var (
	VALIDATE     *[]*Client = &[]*Client{}
	validate     uint32     = uint32(0)
	LOBBY        *[]*Client = &[]*Client{}
	lobby        uint32     = uint32(0)
	GAMES        *[]*Game   = &[]*Game{}
	runningGames uint32     = uint32(0)
	dummyGames   uint32     = uint32(0)
)

type Game struct {
	// TO_ADD     *[2]*[]*websocket.Conn
	TO_ADD     *[2]*[]*Client
	to_add     *[2]uint32
	TO_REPLACE *[2]*[]uint32
	to_replace *[2]uint32
	pro        *[2]uint32
	DAT        *[2]*[]*[20]float64
	PT         *[2]*[]*[34]byte
	CONN       *[2]*[]*websocket.Conn
	NICKS      *[2]*[][20]byte
	exit       bool
}
type Client struct {
	conn   *websocket.Conn
	status uint8
	time   uint
	verify uint8
	nick   [20]byte
}

func main() {
	go LobbyEngine()
	http.HandleFunc("/", HandleEntry)
	fmt.Println("...server started...")
	http.ListenAndServe(":"+GoPort("22222"), nil)
}
func GameEngine() {
	runningGames += 1
	if runningGames > MAX_GAMES {
		fmt.Println("limit games reached")
		runningGames -= 1
		return
	}
	var (
		// TO_ADD     [2]*[]*websocket.Conn = [2]*[]*websocket.Conn{}
		TO_ADD     [2]*[]*Client = [2]*[]*Client{}
		to_add     [2]uint32     = [2]uint32{}
		TO_REPLACE [2]*[]uint32  = [2]*[]uint32{}
		to_replace [2]uint32     = [2]uint32{}
		ET         [2]*[]byte    = [2]*[]byte{}
		PREV       [2]*[]uint32  = [2]*[]uint32{}
		all        [2]uint32     = [2]uint32{}
		pro        [2]uint32     = [2]uint32{}
		DUM        [2]*[]uint32  = [2]*[]uint32{}
		dum        [2]uint32     = [2]uint32{}
		tab        *[]uint32     = &[]uint32{}
		// tab2       *[]*websocket.Conn = &[]*websocket.Conn{}
		tab2  *[]*Client = &[]*Client{}
		pREV  *[]uint32  = &[]uint32{}
		pREV2 *[]uint32  = &[]uint32{}
		// to_ADD     *[]*websocket.Conn = &[]*websocket.Conn{}
		to_ADD     *[]*Client   = &[]*Client{}
		to_REPLACE *[]uint32    = &[]uint32{}
		dUM        *[]uint32    = &[]uint32{}
		eT         *[]byte      = &[]byte{}
		temp       *[10]float64 = &[10]float64{}
		prePt      [20]byte     = [20]byte{3, 0, 0, 0}
		n          uint32
		n2         uint32
		idx        uint32
		idx2       uint32
		CONN       [2]*[]*websocket.Conn
		DAT        [2]*[]*[20]float64
		PT         [2]*[]*[34]byte
		cONN       *[]*websocket.Conn
		dAT        *[]*[20]float64
		pT         *[]*[34]byte
		pT2        *[]*[34]byte
		send       []byte
		pt         *[34]byte
		dat        *[20]float64
		count      uint32
		NICKS      [2]*[][20]byte = [2]*[][20]byte{}
		nICKS      *[][20]byte
	)
	for i := uint32(0); i < 2; i++ {
		to_add[i] = uint32(0)
		to_replace[i] = uint32(0)
		all[i] = uint32(0)
		pro[i] = uint32(0)
		dum[i] = uint32(0)
		// TO_ADD[i] = &[]*websocket.Conn{}
		TO_ADD[i] = &[]*Client{}
		TO_REPLACE[i] = &[]uint32{}
		ET[i] = &[]byte{0}
		PREV[i] = &[]uint32{0}
		DUM[i] = &[]uint32{}
		DAT[i] = &[]*[20]float64{}
		PT[i] = &[]*[34]byte{}
		CONN[i] = &[]*websocket.Conn{}
		NICKS[i] = &[][20]byte{}
	}
	var game *Game = &Game{&TO_ADD, &to_add, &TO_REPLACE, &to_replace, &pro, &DAT, &PT, &CONN, &NICKS, false}
	AddGame(game)
	for {
		for i := uint32(0); i < 2; i++ { //add, replace, calculate, prepare
			to_REPLACE = TO_REPLACE[i]
			dUM = DUM[i]
			eT = ET[i]
			pREV = PREV[i]
			to_ADD = TO_ADD[i]
			dAT = DAT[i]
			pT = PT[i]
			cONN = CONN[i]
			nICKS = NICKS[i]
			if to_replace[i] > 0 { // replace
				*tab = (*tab)[:0]
				to_REPLACE, tab = ChangePointer3(to_REPLACE, tab)
				n = uint32(len(*tab))
				//
				for j := uint32(0); j < n; j++ {
					idx = (*tab)[j]
					// (*cONN)[idx] = nil // ?????????????
					*dUM = append(*dUM, idx)
					(*eT)[idx] = 0
					ReplacedBody(pREV, eT, all[i]-idx, idx)
				}
				dum[i] += n
				to_replace[i] -= n
				pro[i] -= n
				//
				*tab = (*tab)[:0]
				to_REPLACE, tab = ChangePointer3(to_REPLACE, tab)
				*to_REPLACE = append(*to_REPLACE, (*tab)[:]...)
			}
			//
			if to_add[i] > 0 { // add
				*tab2 = (*tab2)[:0]
				to_ADD, tab2 = ChangePointer(to_ADD, tab2)
				n = uint32(len(*tab2))
				if dum[i] >= n {
					for j := uint32(0); j < n; j++ {
						idx = (*dUM)[j]
						(*cONN)[idx] = (*tab2)[j].conn
						(*nICKS)[idx] = (*tab2)[j].nick
						// (*tab2)[j] = nil
						// (*dAT)[idx] = nil // ??????????????????????
						(*dAT)[idx] = InitiateDat(i)
						// SetDat((*dAT)[idx], i)
						(*dAT)[idx][0] = 1
						// Calculate((*dAT)[idx], temp)
						// (*dAT)[idx][0] = 0
						pt = (*pT)[idx]
						pt[2] = byte(idx)
						pt[3] = byte(idx >> 8)
						pt[4] = byte(idx >> 16)
						pt[5] = byte(idx >> 24)
						PreparePt(dat, pt)
						(*eT)[idx] = 1
						ReplacedDummy(pREV, eT, all[i]-idx, idx)
						go HandlePlayer(idx, i, game)
					}
					*dUM = (*dUM)[n:]
					dum[i] -= n
					pro[i] += n
					to_add[i] -= n
				} else {
					if dum[i] > 0 {
						for j := uint32(0); j < dum[i]; j++ {
							idx = (*dUM)[j]
							(*cONN)[idx] = (*tab2)[j].conn
							(*nICKS)[idx] = (*tab2)[j].nick
							// (*tab2)[j] = nil
							// (*dAT)[idx] = nil
							(*dAT)[idx] = InitiateDat(i)
							(*dAT)[idx][0] = 1
							// Calculate((*dAT)[idx], temp)
							// (*dAT)[idx][0] = 0
							// SetDat((*dAT)[idx], i)
							pt = (*pT)[idx]
							pt[2] = byte(idx)
							pt[3] = byte(idx >> 8)
							pt[4] = byte(idx >> 16)
							pt[5] = byte(idx >> 24)
							PreparePt(dat, pt)
							(*eT)[idx] = 1
							ReplacedDummy(pREV, eT, all[i]-idx, idx)
							go HandlePlayer(idx, i, game)
						}
						*dUM = (*dUM)[:0]
						pro[i] += dum[i]
						to_add[i] -= dum[i]
						n -= dum[i]
						dum[i] = 0
					}
					//
					if n > 0 {
						for j := uint32(0); j < n; j++ {
							*cONN = append(*cONN, (*tab2)[j].conn)
							*nICKS = append(*nICKS, (*tab2)[j].nick)
							(*tab2)[j] = nil
							idx = all[i]
							dat = InitiateDat(i)
							dat[0] = 1
							// Calculate(dat, temp)
							// dat[0] = 0
							pt = &[34]byte{0, 0}
							PreparePt(dat, pt)
							*pREV = append(*pREV, idx)
							*eT = append(*eT, 0)
							*dAT = append(*dAT, dat)
							pt[2] = byte(idx)
							pt[3] = byte(idx >> 8)
							pt[4] = byte(idx >> 16)
							pt[5] = byte(idx >> 24)
							*pT = append(*pT, pt)
							(*eT)[all[i]] = 1
							all[i] += 1
							go HandlePlayer(idx, i, game)
						}
						to_add[i] -= n
						pro[i] += n
					}
				}
				*tab2 = (*tab2)[:0]
				to_ADD, tab2 = ChangePointer(to_ADD, tab2)
				*to_ADD = append(*to_ADD, (*tab2)[:]...)
			}
			n = pro[i]
			idx = (*pREV)[all[i]]
			var x uint32
			for j := uint32(0); j < n; j++ { // calculate, prepare
				dat = (*dAT)[idx]
				pt = (*pT)[idx]
				if dat[0] != 0 || dat[1] != 0 {
					Calculate(dat, temp)
					PreparePt(dat, pt)
				}
				x = math.Float32bits(float32(dat[3])) // wysokosc
				pt[30] = byte(x)
				pt[31] = byte(x >> 8)
				pt[32] = byte(x >> 16)
				pt[33] = byte(x >> 24)
				// pt[30] = byte(dat[3])
				idx = (*pREV)[idx]
			}
		}
		cONN = CONN[0]
		pREV = PREV[0]
		pREV2 = PREV[1]
		n = pro[0]
		n2 = pro[1]
		pT = PT[0]
		pT2 = PT[1]
		idx = (*pREV)[all[0]]
		for j := uint32(0); j < n; j++ { // prs A
			prePt[4] = byte(idx)
			prePt[5] = byte(idx >> 8)
			prePt[6] = byte(idx >> 16)
			prePt[7] = byte(idx >> 24)
			idx2 = (*pREV)[all[0]] // t A
			send = send[:0]
			send = append(send, prePt[:]...)
			count = 0
			for k := uint32(0); k < n; k++ { // prs A
				// fmt.Println("sprawdza zasieg a potem kolizje")
				send = append(send, (*pT)[idx2][:]...)
				count += 1
				idx2 = (*pREV)[idx2]
			}
			send[8] = byte(count)
			send[9] = byte(count >> 8)
			send[10] = byte(count >> 16)
			send[11] = byte(count >> 24)
			count = 0
			idx2 = (*pREV2)[all[1]]           // t B
			for k := uint32(0); k < n2; k++ { // prs B
				send = append(send, (*pT2)[idx2][:]...)
				count += 1
				idx2 = (*pREV2)[idx2]
			}
			send[12] = byte(count)
			send[13] = byte(count >> 8)
			send[14] = byte(count >> 16)
			send[15] = byte(count >> 24)
			SendData((*cONN)[idx], &send)
			idx = (*pREV)[idx]
		}
		//
		cONN = CONN[1]
		idx = (*pREV2)[all[1]]
		for j := uint32(0); j < n2; j++ { // prs B
			prePt[4] = byte(idx)              //
			prePt[5] = byte(idx >> 8)         //
			prePt[6] = byte(idx >> 16)        //
			prePt[7] = byte(idx >> 24)        //
			send = send[:0]                   //
			send = append(send, prePt[:]...)  //
			count = 0                         //
			idx2 = (*pREV2)[all[1]]           // t B
			for k := uint32(0); k < n2; k++ { // prs B
				// sprawdzenie warunek dot>a
				send = append(send, (*pT2)[idx2][:]...) //
				count += 1
				idx2 = (*pREV2)[idx2] //
			} //
			send[8] = byte(count)
			send[9] = byte(count >> 8)
			send[10] = byte(count >> 16)
			send[11] = byte(count >> 24)
			count = 0
			idx2 = (*pREV)[all[0]]           // t A
			for k := uint32(0); k < n; k++ { // prs A
				send = append(send, (*pT)[idx2][:]...) //
				count += 1                             //
				idx2 = (*pREV)[idx2]                   //
			}
			send[12] = byte(count)
			send[13] = byte(count >> 8)
			send[14] = byte(count >> 16)
			send[15] = byte(count >> 24)
			SendData((*cONN)[idx], &send)
			idx = (*pREV2)[idx]
		}
		time.Sleep(20 * time.Millisecond)
		if game.exit {
			game = nil
			dummyGames += 1
			runningGames -= 1
			return
		}
	}
}
func LobbyEngine() {
	var (
		TAB  *[]*Client = &[]*Client{}
		NICK *[]*Client = &[]*Client{}
		nick uint32     = uint32(0)
		val  uint32
		n    uint32
	)
	//
	for {
		//
		if validate > 0 {
			*TAB = (*TAB)[:0]
			VALIDATE, TAB = ChangePointer(VALIDATE, TAB)
			TAB, val, n = DoValidate(TAB, NICK)
			validate -= val + n
			nick += val
			VALIDATE, TAB = ChangePointer(VALIDATE, TAB)
			(*VALIDATE) = append((*VALIDATE), (*TAB)[:]...)
			fmt.Println("VALIDATE: ", *VALIDATE)
		}
		//
		if nick > 0 {
			*TAB = (*TAB)[:0]
			NICK, TAB = ChangePointer(NICK, TAB)
			TAB, val, n = DoNick(TAB, LOBBY)
			nick -= val + n
			lobby += val
			NICK, TAB = ChangePointer(NICK, TAB)
			(*NICK) = append((*NICK), (*TAB)[:]...)
			fmt.Println("NICK    : ", *NICK)
		}
		//
		if lobby > 0 {
			*TAB = (*TAB)[:0]
			LOBBY, TAB = ChangePointer(LOBBY, TAB)
			TAB, val, n = DoLobby(TAB)
			lobby -= val + n
			LOBBY, TAB = ChangePointer(LOBBY, TAB)
			(*LOBBY) = append((*LOBBY), (*TAB)[:]...)
			fmt.Println("LOBBY    : ", *LOBBY)
		}
		//
		time.Sleep(time.Second)
	}
}
func ChangePointer(tab1 *[]*Client, tab2 *[]*Client) (*[]*Client, *[]*Client) {
	return tab2, tab1
}
func ChangePointer2(tab1 *[]*websocket.Conn, tab2 *[]*websocket.Conn) (*[]*websocket.Conn, *[]*websocket.Conn) {
	return tab2, tab1
}
func ChangePointer3(tab1 *[]uint32, tab2 *[]uint32) (*[]uint32, *[]uint32) {
	return tab2, tab1
}
func DoValidate(tab *[]*Client, tab2 *[]*Client) (*[]*Client, uint32, uint32) {
	//
	var (
		client  *Client
		err     error
		newTAB  []*Client = []*Client{}
		to_nick uint32    = uint32(0)
		n       uint32    = uint32(0)
	)
	for _, client = range *tab {
		if client.time > 60 {
			client.status = 10
		}
		switch client.status {
		case 0:
			{
				err = client.conn.WriteMessage(2, []byte{0, client.verify})
				if err == nil {
					client.time += 1
					newTAB = append(newTAB, client)
				} else {
					n += 1
					client = nil
				}
			}
		case 1:
			{
				client.time = 0
				*tab2 = append(*tab2, client)
				to_nick += 1
			}
		case 10:
			{
				client.conn.Close()
			}
		default:
			{
				client.conn.Close()
			}

		}
	}
	return &newTAB, to_nick, n
}
func DoNick(tab *[]*Client, tab2 *[]*Client) (*[]*Client, uint32, uint32) {
	//
	var (
		client   *Client
		err      error
		newTAB   []*Client = []*Client{}
		to_lobby uint32    = uint32(0)
		n        uint32    = uint32(0)
	)
	for _, client = range *tab {
		if client.time > 60 {
			client.status = 10
		}
		switch client.status {
		case 1:
			{
				err = client.conn.WriteMessage(2, []byte{1, 0})
				if err == nil {
					client.time += 1
					newTAB = append(newTAB, client)
				} else {
					n += 1
					client = nil
				}
			}
		case 2:
			{
				client.time = 0
				*tab2 = append(*tab2, client)
				to_lobby += 1
			}
		case 10:
			{
				client.conn.Close()
			}
		default:
			{
				client.conn.Close()
			}

		}
	}
	return &newTAB, to_lobby, n
}
func DoLobby(tab *[]*Client) (*[]*Client, uint32, uint32) {
	//
	var (
		client   *Client
		err      error
		newTAB   []*Client = []*Client{}
		to_games uint32    = uint32(0)
		n        uint32    = uint32(0)
		info     []byte
	)
	//
	info = GetInfo(GAMES)
	for _, client = range *tab {
		if client.time > 60 {
			client.status = 10
		}
		switch client.status {
		case 2:
			{
				// fmt.Println("here.............................")
				err = client.conn.WriteMessage(2, info)
				if err == nil {
					client.time += 1
					newTAB = append(newTAB, client)
				} else {
					n += 1
					client = nil
				}
			}
		case 3:
			{
				to_games += 1
				client = nil
			}
		case 10:
			{
				client.conn.Close()
			}
		default:
			{
				client.conn.Close()
			}

		}
	}
	return &newTAB, to_games, n
}
func GoPort(p string) string {
	pt := os.Getenv("PORT")
	if pt == "" {
		pt = p
	}
	return pt
}
func HandleEntry(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/ws/" {
		upgrader := websocket.Upgrader{ReadBufferSize: 128, WriteBufferSize: 1024}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err == nil {
			var client *Client = &Client{conn, 0, 0, 100, [20]byte{0}}
			go HandleClient(client)
			*VALIDATE = append(*VALIDATE, client)
			validate += 1
			return
		} else {
			conn.Close()
			return
		}
	} else {
		fmt.Println(w, "/")
		return
	}
}
func HandleClient(client *Client) {
	//
	var msg []byte
	var err error
	var conn *websocket.Conn = client.conn
	var nick [20]byte = [20]byte{}
	//
	fmt.Println("created client:...", "...", client.nick)
	for {
		_, msg, err = conn.ReadMessage()
		if err == nil {
			fmt.Println("client msg:...", msg)
			switch msg[0] {
			case 0:
				{
					if msg[1] == client.verify-1 {
						client.status = 1
					} else {
						conn.Close()
						return
					}
				}
			case 1:
				{
					if len(msg) == 21 {
						for i := 1; i < 21; i++ {
							nick[i-1] = msg[i]
						}
						client.nick = nick
						client.status = 2
					} else {
						conn.Close()
						return
					}
				}
			case 2:
				{
					if int(msg[1]) > len(*GAMES) {
						conn.Close()
						return
					}
					var game *Game = (*GAMES)[msg[1]]
					if game == nil {
						conn.Close()
						return
					}
					if msg[2] > 1 {
						conn.Close()
						return
					}
					client.status = 3
					// AddConn(conn, game.TO_ADD[msg[2]], &game.to_add[msg[2]])
					AddConn(client, game.TO_ADD[msg[2]], &game.to_add[msg[2]])
					return
				}
			case 3:
				{
					// create game
					go GameEngine()
				}
			}
		} else {
			conn.Close()
			return
		}
	}
}
func AddConn(client *Client, ADD *[]*Client, add *uint32) {
	// func AddConn(conn *websocket.Conn, ADD *[]*websocket.Conn, add *uint32) {
	*ADD = append(*ADD, client)
	*add += 1
}
func AddGame(g *Game) {
	if dummyGames < 1 {
		*GAMES = append(*GAMES, g)
	} else {
		var ix int = -1
		for i, v := range *GAMES {
			if v == nil {
				ix = i
				dummyGames -= 1
				break
			}
		}
		if ix != -1 {
			(*GAMES)[ix] = g
			dummyGames -= 1
		} else {
			*GAMES = append(*GAMES, g)
		}
	}
}
func GetInfo(games *[]*Game) []byte {
	var info []byte = []byte{2, 0}
	var n byte = byte(0)
	var i byte = byte(0)
	for _, game := range *games {
		if game != nil {
			n += 1
			info = append(info, i)
			info = append(info, byte(game.pro[0]))
			info = append(info, byte(game.pro[1]))
		}
		i += 1
	}
	info[1] = n
	return info
}
func HandlePlayer(ix uint32, team uint32, game *Game) {
	// fmt.Println("............handle player..............")
	var (
		msg  []byte
		err  error
		conn *websocket.Conn = (*game.CONN[team])[ix]
		dat  *[20]float64    = (*game.DAT[team])[ix]
	)
	for {
		_, msg, err = conn.ReadMessage()
		if err == nil {
			if msg[0] == 4 {
				dat[0] = float64(msg[3]) - 2
				dat[1] = float64(msg[2]) - 2
				dat[2] = float64(msg[4]) - 2
				// fmt.Println(float64(msg[1]))
				dat[3] += (float64(msg[1]) - dat[3]) * 0.02
				// fmt.Println(dat[3])
			} else if msg[0] == 5 {
				// fmt.Println((*game.NICKS[team])[ix])
				var client *Client = &Client{conn, 2, 0, 100, (*game.NICKS[team])[ix]}
				go HandleClient(client)
				*LOBBY = append(*LOBBY, client)
				lobby += 1
				ReplaceBody(ix, game.TO_REPLACE[team], &game.to_replace[team])
				// fmt.Println("move me to lobby")
				fmt.Println("moved to lobby...........")
				return
			}
		} else {
			ReplaceBody(ix, game.TO_REPLACE[team], &game.to_replace[team])
			conn.Close()
			fmt.Println("end...........HandlePlayer")
			return
		}
	}
}
func ReplaceBody(ix uint32, REPLACE *[]uint32, replace *uint32) {
	*REPLACE = append(*REPLACE, ix)
	*replace += 1
}
func ReplacedBody(prev *[]uint32, et *[]byte, n uint32, ix uint32) {
	var val uint32 = (*prev)[ix]
	var index uint32 = ix
	for i := uint32(0); i < n; i++ {
		index += 1
		if (*et)[index] == 0 {
			(*prev)[index] = val
		} else {
			(*prev)[index] = val
			break
		}
	}
}
func ReplacedDummy(prev *[]uint32, et *[]byte, n uint32, ix uint32) {
	var index uint32 = ix
	for i := uint32(0); i < n; i++ {
		index += 1
		if (*et)[index] == 0 {
			(*prev)[index] = ix
		} else {
			(*prev)[index] = ix
			break
		}
	}
}
func SendData(conn *websocket.Conn, pt *[]byte) {
	// fmt.Println("sending...")
	err := conn.WriteMessage(2, *pt)
	if err != nil {
		fmt.Println("error SendData():", err)
	}
}
func InitiateDat(team uint32) *[20]float64 {
	if team == 0 {
		dat := [20]float64{
			0, 0, 0, //                         ct 0-2
			1,       //                         height 3
			1, 0, 0, //                     	x 4-6
			0, 1, 0, //                      	y 7-9
			0, 0, 1, //                   		z 10-12
			1, 0, 0, 0, //        				q 13-16
			0, 0, 0,
		}
		return &dat
	}
	dat := [20]float64{
		0, 0, 0, //                         ct 0-2
		1,       //                         height
		1, 0, 0, //                     	x 4-6
		0, -1, 0, //                      	y 7-9
		0, 0, -1, //                   		z 10-12
		0, -1, 0, 0, //        				q 13-16
		0, 0, 0,
	}
	return &dat
}
func SetDat(dat *[20]float64, team uint32) {
	if team == 0 {
		dat[0] = 0
		dat[1] = 0
		dat[2] = 0
		dat[3] = 1
		dat[4] = 1
		dat[5] = 0
		dat[6] = 0
		dat[7] = 0
		dat[8] = 1
		dat[9] = 0
		dat[10] = 0
		dat[11] = 0
		dat[12] = 1
		dat[13] = 1
		dat[14] = 0
		dat[15] = 0
		dat[16] = 0
		return
	}
	dat[0] = 0
	dat[1] = 0
	dat[2] = 0
	dat[3] = 1
	dat[4] = 1
	dat[5] = 0
	dat[6] = 0
	dat[7] = 0
	dat[8] = -1
	dat[9] = 0
	dat[10] = 0
	dat[11] = 0
	dat[12] = -1
	dat[13] = 0
	dat[14] = -1
	dat[15] = 0
	dat[16] = 0
}
func Calculate(a *[20]float64, b *[10]float64) {
	b[0] = (a[4]*a[0]*0.001 + a[7]*a[1]*0.01 + a[10]*a[2]*0.001)  // nowy vec obrotu skladowa i
	b[1] = (a[5]*a[0]*0.001 + a[8]*a[1]*0.01 + a[11]*a[2]*0.001)  // nowy vec obrotu skladowa j
	b[2] = (a[6]*a[0]*0.001 + a[9]*a[1]*0.01 + a[12]*a[2]*0.001)  // nowy vec obrotu skladowa k
	b[3] = math.Sqrt(b[0]*b[0] + b[1]*b[1] + b[2]*b[2])           // nowy kat obrotu w rad q0
	b[0] = b[0] / b[3]                                            // skladowa i jednostkowego quat obrotu q1
	b[1] = b[1] / b[3]                                            // skladowa j jednostkowego quat obrotu q2
	b[2] = b[2] / b[3]                                            // skladowa k jednostkowego quat obrotu q3
	b[4] = math.Sin(b[3] / 2)                                     //
	b[5] = math.Cos(b[3] / 2)                                     // q0 od nowego wektora obrotu
	b[6] = b[0] * b[4]                                            // q1 od nowego wektora obrotu
	b[7] = b[1] * b[4]                                            // q2 od nowego wektora obrotu
	b[8] = b[2] * b[4]                                            // q3 od nowego wektora obrotu
	b[0] = b[5]*a[13] - b[6]*a[14] - b[7]*a[15] - b[8]*a[16]      // q0 from Qnew=Qvec*Qprev
	b[1] = b[5]*a[14] + b[6]*a[13] + b[7]*a[16] - b[8]*a[15]      // q1 from Qnew=Qvec*Qprev
	b[2] = b[5]*a[15] + b[7]*a[13] - b[6]*a[16] + b[8]*a[14]      // q2 from Qnew=Qvec*Qprev
	b[3] = b[5]*a[16] + b[8]*a[13] + b[6]*a[15] - b[7]*a[14]      // q3 from Qnew=Qvec*Qprev
	b[4] = 1 / math.Sqrt(b[0]*b[0]+b[1]*b[1]+b[2]*b[2]+b[3]*b[3]) //
	a[13] = b[0] * b[4]                                           // new "summated" quaternion
	a[14] = b[1] * b[4]                                           // new "summated" quaternion
	a[15] = b[2] * b[4]                                           // new "summated" quaternion
	a[16] = b[3] * b[4]                                           // new "summated" quaternion
	a[4] = b[0]*b[0] + b[1]*b[1] - b[2]*b[2] - b[3]*b[3]          // os X
	a[5] = 2 * (b[1]*b[2] + b[0]*b[3])                            //
	a[6] = 2 * (b[1]*b[3] - b[0]*b[2])                            //
	a[7] = 2 * (b[1]*b[2] - b[0]*b[3])                            // os Y
	a[8] = b[0]*b[0] - b[1]*b[1] + b[2]*b[2] - b[3]*b[3]          //
	a[9] = 2 * (b[2]*b[3] + b[0]*b[1])                            //
	a[10] = 2 * (b[0]*b[2] + b[1]*b[3])                           // os Z
	a[11] = 2 * (b[2]*b[3] - b[0]*b[1])                           //
	a[12] = b[0]*b[0] - b[1]*b[1] - b[2]*b[2] + b[3]*b[3]         //
}
func PreCalculate(a *[20]float64, b *[10]float64, rot float64) {
	b[0] = a[4] * rot                                             // + a[7]*dy + a[10]*dz                           //  nowy vec obrotu skladowa i
	b[1] = a[5] * rot                                             // + a[8]*dy + a[11]*dz                           //  nowy vec obrotu skladowa j
	b[2] = a[6] * rot                                             // + a[9]*dy + a[12]*dz                           //  nowy vec obrotu skladowa k
	b[3] = math.Sqrt(b[0]*b[0] + b[1]*b[1] + b[2]*b[2])           // nowy kat obrotu w rad q0
	b[0] = b[0] / b[3]                                            // skladowa i jednostkowego quat obrotu q1
	b[1] = b[1] / b[3]                                            // skladowa j jednostkowego quat obrotu q2
	b[2] = b[2] / b[3]                                            // skladowa k jednostkowego quat obrotu q3
	b[4] = math.Sin(b[3] / 2)                                     //
	b[5] = math.Cos(b[3] / 2)                                     // q0 od nowego wektora obrotu
	b[6] = b[0] * b[4]                                            // q1 od nowego wektora obrotu
	b[7] = b[1] * b[4]                                            // q2 od nowego wektora obrotu
	b[8] = b[2] * b[4]                                            // q3 od nowego wektora obrotu
	b[0] = b[5]*a[13] - b[6]*a[14] - b[7]*a[15] - b[8]*a[16]      // q0 from Qnew=Qvec*Qprev
	b[1] = b[5]*a[14] + b[6]*a[13] + b[7]*a[16] - b[8]*a[15]      // q1 from Qnew=Qvec*Qprev
	b[2] = b[5]*a[15] + b[7]*a[13] - b[6]*a[16] + b[8]*a[14]      // q2 from Qnew=Qvec*Qprev
	b[3] = b[5]*a[16] + b[8]*a[13] + b[6]*a[15] - b[7]*a[14]      // q3 from Qnew=Qvec*Qprev
	b[4] = 1 / math.Sqrt(b[0]*b[0]+b[1]*b[1]+b[2]*b[2]+b[3]*b[3]) //
	a[13] = b[0] * b[4]                                           // new "summated" quaternion
	a[14] = b[1] * b[4]                                           // new "summated" quaternion
	a[15] = b[2] * b[4]                                           // new "summated" quaternion
	a[16] = b[3] * b[4]                                           // new "summated" quaternion
	a[4] = b[0]*b[0] + b[1]*b[1] - b[2]*b[2] - b[3]*b[3]          // os X
	a[5] = 2 * (b[1]*b[2] + b[0]*b[3])                            //
	a[6] = 2 * (b[1]*b[3] - b[0]*b[2])                            //
	a[7] = 2 * (b[1]*b[2] - b[0]*b[3])                            // os Y
	a[8] = b[0]*b[0] - b[1]*b[1] + b[2]*b[2] - b[3]*b[3]          //
	a[9] = 2 * (b[2]*b[3] + b[0]*b[1])                            //
	a[10] = 2 * (b[0]*b[2] + b[1]*b[3])                           // os Z
	a[11] = 2 * (b[2]*b[3] - b[0]*b[1])                           //
	a[12] = b[0]*b[0] - b[1]*b[1] - b[2]*b[2] + b[3]*b[3]         //
}
func CalculateCol(a *[20]float64, b *[12]float64, c *[3]float64) {
	b[0] = -c[0] * 0.1
	b[1] = -c[1] * 0.1
	b[2] = -c[2] * 0.1
	b[3] = math.Sqrt(b[0]*b[0] + b[1]*b[1] + b[2]*b[2])           // nowy kat obrotu w rad q0
	b[0] = b[0] / b[3]                                            //                          skladowa i jednostkowego quat obrotu q1
	b[1] = b[1] / b[3]                                            //                          skladowa j jednostkowego quat obrotu q2
	b[2] = b[2] / b[3]                                            //                          skladowa k jednostkowego quat obrotu q3
	b[4] = math.Sin(b[3] / 2)                                     //
	b[5] = math.Cos(b[3] / 2)                                     //							    q0 od nowego wektora obrotu
	b[6] = b[0] * b[4]                                            //                          q1 od nowego wektora obrotu
	b[7] = b[1] * b[4]                                            //                          q2 od nowego wektora obrotu
	b[8] = b[2] * b[4]                                            //                          q3 od nowego wektora obrotu
	b[0] = b[5]*a[13] - b[6]*a[14] - b[7]*a[15] - b[8]*a[16]      // q0 from Qnew=Qvec*Qprev
	b[1] = b[5]*a[14] + b[6]*a[13] + b[7]*a[16] - b[8]*a[15]      // q1 from Qnew=Qvec*Qprev
	b[2] = b[5]*a[15] + b[7]*a[13] - b[6]*a[16] + b[8]*a[14]      // q2 from Qnew=Qvec*Qprev
	b[3] = b[5]*a[16] + b[8]*a[13] + b[6]*a[15] - b[7]*a[14]      // q3 from Qnew=Qvec*Qprev
	b[4] = 1 / math.Sqrt(b[0]*b[0]+b[1]*b[1]+b[2]*b[2]+b[3]*b[3]) //
	a[13] = b[0] * b[4]                                           // new "summated" quaternion
	a[14] = b[1] * b[4]                                           // new "summated" quaternion
	a[15] = b[2] * b[4]                                           // new "summated" quaternion
	a[16] = b[3] * b[4]                                           // new "summated" quaternion
	a[4] = b[0]*b[0] + b[1]*b[1] - b[2]*b[2] - b[3]*b[3]
	a[5] = 2 * (b[1]*b[2] + b[0]*b[3])
	a[6] = 2 * (b[1]*b[3] - b[0]*b[2])
	a[7] = 2 * (b[1]*b[2] - b[0]*b[3])
	a[8] = b[0]*b[0] - b[1]*b[1] + b[2]*b[2] - b[3]*b[3]
	a[9] = 2 * (b[2]*b[3] + b[0]*b[1])
	a[10] = 2 * (b[0]*b[2] + b[1]*b[3])
	a[11] = 2 * (b[2]*b[3] - b[0]*b[1])
	a[12] = b[0]*b[0] - b[1]*b[1] - b[2]*b[2] + b[3]*b[3]
}
func PreparePt(dat *[20]float64, pt *[34]byte) {
	//
	x := math.Float32bits(float32(dat[4]))
	pt[6] = byte(x)
	pt[7] = byte(x >> 8)
	pt[8] = byte(x >> 16)
	pt[9] = byte(x >> 24)
	//
	x = math.Float32bits(float32(dat[5]))
	pt[10] = byte(x)
	pt[11] = byte(x >> 8)
	pt[12] = byte(x >> 16)
	pt[13] = byte(x >> 24)
	//
	x = math.Float32bits(float32(dat[6]))
	pt[14] = byte(x)
	pt[15] = byte(x >> 8)
	pt[16] = byte(x >> 16)
	pt[17] = byte(x >> 24)
	//
	x = math.Float32bits(float32(dat[7]))
	pt[18] = byte(x)
	pt[19] = byte(x >> 8)
	pt[20] = byte(x >> 16)
	pt[21] = byte(x >> 24)
	//
	x = math.Float32bits(float32(dat[8]))
	pt[22] = byte(x)
	pt[23] = byte(x >> 8)
	pt[24] = byte(x >> 16)
	pt[25] = byte(x >> 24)
	//
	// x = math.Float32bits(float32(dat[9]))
	// pt[26] = byte(x)
	// pt[27] = byte(x >> 8)
	// pt[28] = byte(x >> 16)
	// pt[29] = byte(x >> 24)
	//
	// x = math.Float32bits(float32(dat[3])) // wysokosc
	// pt[30] = byte(x)
	// pt[31] = byte(x >> 8)
	// pt[32] = byte(x >> 16)
	// pt[33] = byte(x >> 24)
	// return
	x = math.Float32bits(float32(dat[9]))
	pt[26] = byte(x)
	pt[27] = byte(x >> 8)
	pt[28] = byte(x >> 16)
	pt[29] = byte(x >> 24)
}
