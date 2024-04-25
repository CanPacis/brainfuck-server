package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/CanPacis/bfx-server/bfx"
)

func main() {
	socket, err := create_socket()
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	for {
		w, e := socket.Accept()
		if e != nil {
			log.Println(e.Error())
		}

		header, _ := os.OpenFile("./header.txt", os.O_RDONLY, 0644)
		content, _ := os.OpenFile("./content.html", os.O_RDONLY, 0644)
		bfx.RunFile("server.bf", bfx.Cell(header.Fd()), bfx.Cell(content.Fd()), bfx.Cell(w.fd))
		time.Sleep(time.Millisecond * 100)
		header.Close()
		content.Close()
		w.Close()
	}
}

func create_socket() (*netSocket, error) {
	ipFlag := flag.String("ip_addr", "127.0.0.1", "The IP address to use")
	portFlag := flag.Int("port", 80, "The port to use.")
	flag.Parse()

	ip := net.ParseIP(*ipFlag)
	port := *portFlag
	socket, err := newNetSocket(ip, port)
	if err != nil {
		return nil, err
	}

	return socket, nil
}
