package main

import (
	"./protocol"
	"fmt"
	"net"
	"os"
)

func main() {
	netListen, err := net.Listen("tcp", ":9988")
	CheckError(err)
	defer netListen.Close()
	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
		go senderToClient(conn)

	}
}

func senderToClient(conn net.Conn) {
	conn.Write(protocol.Packet([]byte("欢迎光临")))
}

func handleConnection(conn net.Conn) {
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)
	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
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
