package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Apply() {
	fmt.Println("Apply")
}

func Read() {
	fmt.Println("Read")
}

func Restore() {
	fmt.Println("Restore")
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
	router := mux.NewRouter()
	router.HandleFunc("/{key}", get).Methods("GET")
	router.HandleFunc("/", set).Methods("PUT")
	router.HandleFunc("/nodes/list", listNodes).Methods("GET")
	router.HandleFunc("/nodes/{id}", removeNode).Methods("DELETE")

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
