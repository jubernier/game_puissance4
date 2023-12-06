package main

import (
	"log"
	"net"
	"sync"
	"strconv"
)

var msgChannel = make(chan string)
var client1Ready = make(chan struct{})
var client2Ready = make(chan struct{})
var wg sync.WaitGroup
var logLock sync.Mutex

func handleClient(c net.Conn, clientID int) {
	switch clientID {
	case 1:
		close(client1Ready)
	case 2:
		close(client2Ready)
	}

	message := "Hello, Client " + strconv.Itoa(clientID)
	_, err := c.Write([]byte(message))
	if err != nil {
		log.Println("Write error:", err)
		return
	}
	log.Println("Server sent:", message)
}

func handleClients(c1 net.Conn, c2 net.Conn) {
	message := "Tous les clients sont là"
	logLock.Lock()
	_, err := c1.Write([]byte(message))
	if err != nil {
		log.Println("Write error:", err)
		return
	}
	logLock.Unlock()
	logLock.Lock()
	_, err = c2.Write([]byte(message))
	if err != nil {
		log.Println("Write error:", err)
		return
	}
	logLock.Unlock()
	log.Println("Server sent:", message)
}

func main() {
	var list_connection []net.Conn

	// synchroniser se qui se passe entre les deux clients
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	defer listener.Close()
	log.Println("Server listening on :8080")
	wg.Add(2)

	// Accept clients concurrently
	for i := 0; i < 2; i++ {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		list_connection = append(list_connection, conn)

		go func(c net.Conn, id int) {
			defer wg.Done()
			handleClient(c, id)
		}(conn, i+1)
	}

	// Wait for both clients to be ready
	wg.Wait()

	// Handle clients after they are both ready
	if len(list_connection) == 2 {
		handleClients(list_connection[0], list_connection[1])
	}
	// Ensure all connections are closed when main exits
	defer func() {
		for _, conn := range list_connection {
			conn.Close()
		}
	}()


	<-client1Ready
	<-client2Ready
}
