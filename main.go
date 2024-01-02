package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/memberlist"
)

type Broadcast struct {
	message []byte
}

func (b *Broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *Broadcast) Message() []byte {
	return b.message
}

func (b *Broadcast) Finished() {
	log.Printf("Finished broadcasting %s", string(b.message))
}

func newBroadcast(message []byte) *Broadcast {
	return &Broadcast{message: message}
}

type Server struct {
	list       *memberlist.Memberlist
	broadcasts *memberlist.TransmitLimitedQueue
}

func (s *Server) NotifyJoin(node *memberlist.Node) {
	log.Printf("Node %s joined", node.Name)
}

func (s *Server) NotifyLeave(node *memberlist.Node) {
	log.Printf("Node %s left", node.Name)
}

func (s *Server) NotifyUpdate(node *memberlist.Node) {
	log.Printf("Node %s updated", node.Name)
}

func (s *Server) NodeMeta(limit int) []byte {
	return []byte{}
}

func (s *Server) NotifyMsg(b []byte) {
	log.Printf("Received message: %s", string(b))
}

func (s *Server) GetBroadcasts(overhead, limit int) [][]byte {
	return s.broadcasts.GetBroadcasts(overhead, limit)
}

func (s *Server) LocalState(join bool) []byte {
	return []byte{}
}

func (s *Server) MergeRemoteState(buf []byte, join bool) {}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}

		s.broadcasts.QueueBroadcast(newBroadcast(body))
		w.Write([]byte("Message broadcasted"))
	}

	if r.Method == http.MethodGet {
		s.ListMembers(w, r)
	}

}

func (s *Server) ListMembers(w http.ResponseWriter, r *http.Request) {
	members := s.list.Members()
	for _, member := range members {
		fmt.Fprintf(w, "%s\n", member.Name)
	}
}

func main() {

	portP := flag.Int("port", 0, "port to listen on")
	listenP := flag.String("listen", ":8080", "address to listen on")
	nameP := flag.String("name", "node", "name of the node")
	peersP := flag.String("peers", "", "comma separated list of peers")
	flag.Parse()

	server := &Server{}

	config := memberlist.DefaultLocalConfig()
	config.Name = *nameP
	config.BindPort = *portP
	config.Events = server
	config.Delegate = server

	logFile, _ := os.OpenFile(fmt.Sprintf("node-%d.log", *portP), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(logFile)

	var err error
	server.list, err = memberlist.Create(config)
	if err != nil {
		log.Fatal(err)
	}

	server.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes:       func() int { return server.list.NumMembers() },
		RetransmitMult: 3,
	}

	if *peersP != "" {
		peers := strings.Split(*peersP, ",")
		log.Printf("Joining %d peers", len(peers))
		_, err = server.list.Join(peers)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Listening on %s", *listenP)
	log.Fatal(http.ListenAndServe(*listenP, server))

}
