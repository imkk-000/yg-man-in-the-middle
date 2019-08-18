package main

import (
	"bufio"
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

	switch {
	case (data[0] == 0x64 && data[1] == 0x80):
		game = true
		b := []byte{}
		return len(b), b
	}
	return n, data[:n]
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
