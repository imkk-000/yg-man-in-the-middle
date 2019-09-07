package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"net"
	"yulgang/model"
)

const BUFFERSIZE = 2048

var ServerConfig model.Config

func WriteLogFile(data []byte, prefix string) error {
	f, err := ioutil.TempFile("log", prefix)
	f.Write(data)
	f.Close()
	return err
}

func InjectData(n int, data []byte) (int, []byte) {
	if n <= 0 {
		return n, data
	}
	code := binary.LittleEndian.Uint16(data[:2])
	switch code {
	case 0x8064:
		data = data[4:n]
		// get ip
		ipSize := binary.LittleEndian.Uint16(data)
		ipAddr := string(data[2 : ipSize+2])
		port := binary.LittleEndian.Uint16(data[ipSize+2:])
		data = data[ipSize+2+2:]
		userSize := binary.LittleEndian.Uint16(data)
		user := string(data[2 : userSize+2])
		log.Println(ipAddr, port, userSize, user, data)

		// new data
		newIPAddr := ServerConfig.Game.IP
		newPort := ServerConfig.Game.Port

		newIPAddrLenBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(newIPAddrLenBytes, uint16(len(newIPAddr)))
		newPortBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(newPortBytes, uint16(newPort))

		// create new packet data
		newData := append(newIPAddrLenBytes, []byte(newIPAddr)...)
		newData = append(newData, newPortBytes...)
		newData = append(newData, data...)
		newData = BuilderNewPacket(code, newData)
		return len(newData), newData
	}
	return n, data[:n]
}

func BuilderNewPacket(code uint16, data []byte) (newPacket []byte) {
	newPacket = make([]byte, 2+2)
	binary.LittleEndian.PutUint16(newPacket, code)
	binary.LittleEndian.PutUint16(newPacket[2:], uint16(len(data)))
	newPacket = append(newPacket, data...)
	return
}

func WriteData(writer io.Writer, reader io.Reader) (data []byte, n int, err error) {
	data = make([]byte, BUFFERSIZE)
	n, err = reader.Read(data)
	n, data = InjectData(n, data)
	writer.Write(data)
	return
}

func handler(config model.Config, c *net.TCPConn) {
	defer c.Close()
	loginNetConfig := &net.TCPAddr{
		IP:   net.ParseIP(config.Login.IP),
		Port: config.Login.Port,
	}
	loginClient, err := net.DialTCP("tcp", nil, loginNetConfig)
	if err != nil {
		log.Print(err)
	}
	defer loginClient.Close()

	// reuse err
	err = nil

	for {
		_, n, err := WriteData(loginClient, c)
		log.Println("Send", n, err)
		if err == io.EOF {
			break
		}

		_, n, err = WriteData(c, loginClient)
		log.Println("Recv", n, err)
		if err == io.EOF {
			break
		}
	}
}

func ServerListener(config model.Config) {
	serverNetConfig := &net.TCPAddr{
		IP:   net.ParseIP(config.Server.IP),
		Port: config.Server.Port,
	}
	server, err := net.ListenTCP("tcp", serverNetConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server listen on port", config.Server.Port)

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("client connected on", conn.RemoteAddr().String())

		go handler(config, conn)
	}
}

func main() {
	ServerConfig = ServerConfig.Read("config.json")
	if !ServerConfig.Log {
		emptyBuffer := bufio.NewWriter(nil)
		log.SetOutput(emptyBuffer)
	}
	ServerListener(ServerConfig)
}
