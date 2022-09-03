package server

import (
	"log"
	"net"
)

type Handler = func([] byte) []byte


type Server struct {
	port uint16
	bufferSize uint16
	handler Handler
}

func NewServer(port uint16, bufferSize uint16, handler Handler) Server {
	server := Server{port: port, bufferSize: bufferSize, handler: handler}
	return server
}


func (server *Server) Run() {
	log.Printf("Server start listening port: %d", server.port)
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP: net.IPv4(0, 0, 0, 0),
		Port: int(server.port),
	})
	
	if err != nil {
		log.Fatalf("net.ListenUDP failed - %s", err)
	}
	
	defer socket.Close()
	
	data := make([]byte, server.bufferSize)
	
	for {
		log.Printf("Waiing for connections ...")
		read, remoteAddress, err := socket.ReadFromUDP(data)
		if err != nil {
			log.Printf("socket.ReadFromUDP failed - %s", err)
			continue
		}
		
		log.Printf("Client connected %s. Read %d bytes.", remoteAddress, read)
		
		response := server.handler(data)
		_, err = socket.WriteToUDP(response, remoteAddress)
		if err != nil {
			log.Printf("socket.WriteToUDP failed - %s", err)
		}
	}
}
