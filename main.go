package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/raft"
)

type kvFsm struct {
	v map[string]string
}

type snapshotNoop struct{}

func (sn snapshotNoop) Persist(_ raft.SnapshotSink) error { return nil }
func (sn snapshotNoop) Release()

func (kf *kvFsm) Apply(log *raft.Log) any {
	fmt.Println("Apply")
	return true
}

func (kf *kvFsm) Restore(io.ReadCloser) error {
	fmt.Println("Restore")
	return nil
}

func (kf *kvFsm) Snapshot() (raft.FSMSnapshot, error) {
	return snapshotNoop{}, nil
}

func setupRaft(dir string, nodeId int, raftAddress string, kvFsm *kvFsm) (*raft.Raft, error) {
	snapshots, err := raft.NewFileSnapshotStore(path.Join(dir, "snapshot"), 2, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("Could not create snapshot store: %s", err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", raftAddress)
	if err != nil {
		return nil, fmt.Errorf("Could not resolve address: %s", err)
	}

	transport, err := raft.NewTCPTransport(raftAddress, tcpAddr, 10, time.Second*10, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("Could not create tcp transport: %s", err)
	}

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeId)

	r, err := raft.NewRaft(config, kvFsm, logs, store, snapshots, transport)
	if err != nil {
		return nil, fmt.Errorf("Could not create raft instance: %s", err)
	}

	return r, nil
}

func getConfig() {
}

func get(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	fmt.Println(key)
}

func set(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Set key value pair")
}

func listNodes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List nodes")
}

func removeNode(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
}

func main() {
	dataDir := "data"

	r, err := setupRaft(path.Join(dataDir, "raft"+config.id), config.id, "localhost:9090", kf)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/{key}", get).Methods("GET")
	router.HandleFunc("/", set).Methods("PUT")
	router.HandleFunc("/nodes/list", listNodes).Methods("GET")
	router.HandleFunc("/nodes/{id}", removeNode).Methods("DELETE")

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
