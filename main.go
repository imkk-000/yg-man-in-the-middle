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

func GetData8064(data []byte, newIPAddr string, newPort int) ([]byte, model.IpConfig, string) {
	// get ip & port & username
	ipSize := binary.LittleEndian.Uint16(data)
	gameConfig := model.IpConfig{
		IP:   string(data[2 : ipSize+2]),
		Port: int(binary.LittleEndian.Uint16(data[ipSize+2:])),
	}
	data = data[ipSize+2+2:]
	userSize := binary.LittleEndian.Uint16(data)
	user := string(data[2 : userSize+2])

	// prepare new data
	newIPAddrLenBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(newIPAddrLenBytes, uint16(len(newIPAddr)))
	newPortBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(newPortBytes, uint16(newPort))

	// create new packet data
	newData := append(newIPAddrLenBytes, []byte(newIPAddr)...)
	newData = append(newData, newPortBytes...)
	newData = append(newData, data...)

	log.Println("created:", gameConfig, "to", ServerConfig.Game, user)
	return newData, gameConfig, user
}

func InjectData(n int, data []byte) (int, []byte) {
	if n <= 0 {
		return n, data
	}
	code := binary.LittleEndian.Uint16(data[:2])
	switch code {
	case 0x8064:
		data, _, _ = GetData8064(data[4:n], ServerConfig.Game.IP, ServerConfig.Game.Port)
		data = BuilderNewPacket(code, data)
		return len(data), data
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

func handler(c *net.TCPConn) {
	defer c.Close()
	loginNetConfig := &net.TCPAddr{
		IP:   net.ParseIP(ServerConfig.Login.IP),
		Port: ServerConfig.Login.Port,
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

func ServerListener() {
	serverNetConfig := &net.TCPAddr{
		IP:   net.ParseIP(ServerConfig.Server.IP),
		Port: ServerConfig.Server.Port,
	}
	server, err := net.ListenTCP("tcp", serverNetConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server listen on port", ServerConfig.Server.Port)

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("client connected on", conn.RemoteAddr().String())

		go handler(conn)
	}
}

func main() {
	ServerConfig.Read("config.json")
	if !ServerConfig.Log {
		emptyBuffer := bufio.NewWriter(nil)
		log.SetOutput(emptyBuffer)
	}
	ServerListener()
}
