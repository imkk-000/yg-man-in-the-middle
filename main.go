package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"yulgang/model"
)

const BUFFERSIZE = 2048

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
	code := binary.BigEndian.Uint16(data[:2])
	switch code {
	case 0x6480:
		fmt.Println(data)
	}
	return n, data[:n]
}

func BuilderNewPacket(code uint16, data []byte) (newPacket []byte) {
	newPacket = make([]byte, 2)
	binary.BigEndian.PutUint16(newPacket, code)
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, uint16(len(data)))
	newPacket = append(newPacket, lenBytes...)
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
	serverConfig := model.Config{}
	serverConfig = serverConfig.Read("config.json")
	if !serverConfig.Log {
		emptyBuffer := bufio.NewWriter(nil)
		log.SetOutput(emptyBuffer)
	}
	ServerListener(serverConfig)
}
