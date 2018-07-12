package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	State *State `json:"state,omitempty"`
}

type State struct {
	GamesPlayed int `json:"gamesplayed,omitempty"`
	Score       int `json:"score,omitempty"`
}

var users []User

func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users[0])
}
func GetState(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users[0].State)
}
func GetFriends(w http.ResponseWriter, r *http.Request) {
	var friends []User
	payload := map[string]interface{}{
		"friends": friends,
	}
	json.NewEncoder(w).Encode(payload)
}
func GetAll(w http.ResponseWriter, r *http.Request) {
	payload := map[string]interface{}{
		"users": users,
	}
	json.NewEncoder(w).Encode(payload)
}

// our main function
func main() {

	users = append(users, User{ID: "18dd75e9-3d4a-48e2-bafc-3c8f95a8f0d1", Name: "jj", State: &State{GamesPlayed: 42, Score: 358}})
	users = append(users, User{ID: "3123123131-3d4a-48e2-bafc-3c8f95a8f0d1", Name: "kk", State: &State{GamesPlayed: 412, Score: 3123}})

	router := mux.NewRouter()
	router.HandleFunc("/user", CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}/state", GetState).Methods("PUT", "GET")
	router.HandleFunc("/user/{id}/friends", GetState).Methods("PUT", "GET")
	router.HandleFunc("/user", GetAll).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
