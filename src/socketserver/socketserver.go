package socketserver

import (
	"os"
	"fmt"
	"net"
	"../github.com/sluu99/uuid"
	"../protocol"
)


var sessions map[string]net.Conn = make(map[string]net.Conn)

func StartSocket(servernetwork string,serveraddress string, flag chan bool){
	netListen, err := net.Listen(servernetwork, serveraddress)
	CheckError(err)
	defer func(){
		netListen.Close()
		flag <- true
	}()
	Log("socket server start success on ", serveraddress,servernetwork)
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		//加入到session中
		addSession(conn)
		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
		go senderToClient(conn,[]byte("welcome"))
	}
}

func addSession(conn net.Conn){
	id := uuid.Rand()
	sessions[id.Hex()] = conn
}

func removeSession(conn net.Conn){
	for key ,val := range sessions{
		if val == conn{
			delete(sessions,key)
		}
	}
}

func Broadcast(msg []byte){
	for _,val := range sessions{
		senderToClient(val,msg)
	}
}

func senderToClient(conn net.Conn, msg []byte) {
	conn.Write(protocol.Packet(msg))
}

func handleConnection(conn net.Conn) {
	tmpBuffer := make([]byte, 0)
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), "client connect error!", err)
			removeSession(conn)
			return
		}
		tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			Log(string(data))
		}
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

