package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"yulgang/model"
)

const BUFFERSIZE = 2048

func WriteData(writer io.Writer, reader io.Reader) (n int, err error) {
	data := make([]byte, BUFFERSIZE)
	n, err = reader.Read(data)
	data = data[:n]
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

	for {
		n, err := WriteData(loginClient, c)
		log.Println("Send", n, err)
		if err == io.EOF {
			break
		}

		n, err = WriteData(c, loginClient)
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
