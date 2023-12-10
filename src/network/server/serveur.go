package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"network"
	"sync"
)

// Server : Structure de données avec tout ce qu'il nous faut !
type Server struct {
	clients  []net.Conn // Tableau de connexions qui représentent les clients
	listener net.Listener

	writeChans []chan string // Tunnel d'écriture
	readChans  []chan string // Tunnel de lecture

	countTokens     int        // Compteur utilisé plusieurs fois dans le programme
	mutexCountRC    sync.Mutex // Mutex permettant de bloquer certaines parties du code
	finishTimes     [2]string  // Tableau des temps
	runnerPositions [2]int     // Tableau des runners sélectionnés
}

// newServer : Creating a new server with clients, listener, writeChans, readChans
func newServer() Server {
	listener, _ := net.Listen("tcp", ":8080")
	return Server{
		[]net.Conn{},
		listener,
		[]chan string{},
		[]chan string{},
		0,
		sync.Mutex{},
		[2]string{},
		[2]int{0, 1},
	}
}

// acceptClients : Accepte 2 connexions de clients
func (s *Server) acceptClients() {
	for i := 0; i < 2; i++ {
		conn, err := s.listener.Accept() // Accepte la connection du client
		if err != nil {                  // La méthode Accept est bloquante, elle ne passe à l'étape suivante (ligne suivante) que lorsqu'une connection vient d'être acceptée
			log.Fatal(err)
		}
		s.clients = append(s.clients, conn) // Ajout de la connection de ce nouveau client à la liste des clients

		// Créer les canaux de communication avec ce client
		// Enregistre ces canaux dans un tableau répertoriant tous les canaux
		s.writeChans = append(s.writeChans, make(chan string, 10))
		s.readChans = append(s.readChans, make(chan string, 10))

		// Initialiser les goroutines de communication de ce client avec la connection et le channel initialisé précédement
		go ReadFromNetWork(bufio.NewReader(conn), s.readChans[i])
		go WriteFromNetWork(bufio.NewWriter(conn), s.writeChans[i])

		s.writeChans[i] <- network.CLIENT_NUMBER + fmt.Sprint(i) // Envoyer l'ID au client qui vient de se connecter
		s.sendToAll(network.CLIENTS_IN_QUEUE + fmt.Sprint(i+1))  // Envoie à tout le monde le nombre de clients actuellement connectés

		log.Println("Client ", i, "connected")
		log.Println("in queue: ", fmt.Sprint(i+1))
	}
}

// sendToAll : Envoyer un message à tous les clients
func (s *Server) sendToAll(message string) {
	for _, chann := range s.writeChans { // Pour tous les channels d'écriture ouverts
		chann <- message // On envoie le message dans ce canal
	}
}

// ReadFromNetWork Fonction permettant de lire en boucle sur un reader
func ReadFromNetWork(connReader *bufio.Reader, channel chan string) {
	for {
		msg, err := connReader.ReadString('\n')
		log.Print(msg)
		if err != nil {
			log.Panic(err)
		}
		channel <- msg[:len(msg)-1] // Ajoute le message au channel en retirant le dernier caractère (delimiter)
	}
}

func WriteFromNetWork(connWriter *bufio.Writer, channel chan string) {
	var message string
	for {
		message = <-channel
		_, err := connWriter.WriteString(message + "\n")
		if err != nil {
			log.Panic(err)
		}
		connWriter.Flush()
	}
}

// comClient : Communication avec un client
func (s *Server) comClient(clientID int) {
	var message string
	for {
		message = <-s.readChans[clientID] // On lit le message du channel du client
		// On rappelle que cette fonction est lancée comme go routine et que le serveur lance cette fonction pour chaque client
		switch message[:1] {
		// Cas lorsqu'un client vient de se connecter
		case network.CLIENT_CONNECTED:
			log.Println("client ", clientID, "connecté")

		// Cas lorsqu'un client vient de sélectionner un jeton
		case network.CLIENT_CHOOSE_TOKEN:
			var selected = message[1:]
			s.mutexCountRC.Lock() // On bloque la partie du code suivante
			if selected == "0" {  // Le client a déselectionné son jeton
				s.countTokens--
				log.Println("client", clientID, "has deselected his runner")
			}
			if selected == "1" { // Le client a sélectionné son jeton
				s.countTokens++ // On compte le nombre de joueurs ayant choisi leur jeton
				log.Println("client", clientID, "chose his runner")
			}
			s.sendToAll(network.CLIENTS_IN_QUEUE + fmt.Sprint(s.countTokens)) // On envoie à tous les clients le nombre de joueurs ayant choisi leur personnage
			if s.countTokens == 2 {                                           // Si les 2 joueurs ont choisi leur personnage
				s.sendToAll(network.ALL_TOKEN_CHOOSEN) // On envoie à tous les clients que tout le monde a choisi son personnage
				s.sendToAll(network.START_RACE)        // On envoie à tous les clients que la course commence

				s.countTokens = 0 // On réinitialise notre compteur
			}
			s.mutexCountRC.Unlock() // On débloque cette partie du code
		}
	}
}

func main() {
	var server = newServer() // Création du serveur
	log.Println("Listening for connections")
	server.acceptClients() // Attend que les 2 clients soient connectés
	// Le programme ne passe pas à l'étape suivante (ligne suivante) tant que 4 clients ne se sont pas connectés
	server.sendToAll(network.ALL_CONNECTED) // Envoie à tous les clients l'information que tous les clients sont connectés
	log.Println("All clients connected")
	defer server.listener.Close()         // Fermer la connection à la fin du programme
	for i, conn := range server.clients { // Pour tous les clients connectés
		defer conn.Close()
		go server.comClient(i) // On initialise une connection avec ce client
	}

	for {
	} // Pour empêcher le serveur de s'éteindre
}
