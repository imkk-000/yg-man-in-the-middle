package main

import (
	"io"
	"log"
	"net"
	"yulgang/model"
)

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

	for {
		data := make([]byte, 2048)
		n, err := c.Read(data)
		log.Println("send", n, err)
		data = data[:n]
		loginClient.Write(data)
		if err == io.EOF {
			break
		}

		data = make([]byte, 2048)
		n, err = loginClient.Read(data)
		log.Println("receive", n, err)
		data = data[:n]
		c.Write(data)
		if err == io.EOF {
			break
		}
	}
}

func main() {
	serverConfig := model.Config{}
	serverConfig = serverConfig.Read("config.json")

	serverNetConfig := &net.TCPAddr{
		IP:   net.ParseIP(serverConfig.Server.IP),
		Port: serverConfig.Server.Port,
	}
	server, err := net.ListenTCP("tcp", serverNetConfig)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}

		go handler(serverConfig, conn)
	}
}
