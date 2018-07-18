package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var dal *Dal

type Player struct {
	ID          uuid.UUID `db:"id,omitempty" json:"id,omitempty"`
	Name        string    `db:"name,omitempty" json:"name,omitempty"`
	GamesPlayed int       `db:"games" json:"gamesplayed,omitempty"`
	Score       int       `db:"score" json:"score,omitempty"`
}

type Friend struct {
	ID    uuid.UUID `db:"id,omitempty" json:"id,omitempty"`
	Name  string    `db:"name,omitempty" json:"name"`
	Score int       `db:"score" json:"highscore"`
}

type FriendsRequest struct {
	IDs *[]uuid.UUID `json:"friends"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	dal, _ := dal.GetSession()

	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 400)
		return
	}

	newPlayer := Player{
		ID:   uuid.New(),
		Name: player.Name,
	}

	err = dal.InsertPlayer(newPlayer)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("user created")
	json.NewEncoder(w).Encode(newPlayer)
}

func SaveState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	dal, _ := dal.GetSession()

	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = dal.UpdateState(player.Score, player.GamesPlayed, id)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("State updated")
}

func GetState(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	dal, _ := dal.GetSession()

	var player Player
	err := dal.GetState(&player, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	payload := map[string]interface{}{
		"gamesPlayed": player.GamesPlayed,
		"score":       player.Score,
	}
	json.NewEncoder(w).Encode(payload)
}

func GetFriends(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]
	dal, _ := dal.GetSession()

	friends, err := dal.GetFriends(id)

	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	payload := map[string]interface{}{
		"friends": friends,
	}
	log.Print("Friendlist retrieved with success")
	json.NewEncoder(w).Encode(payload)
}

func AddFriends(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	dal, _ := dal.GetSession()

	var friends FriendsRequest
	err := json.NewDecoder(r.Body).Decode(&friends)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = dal.UpdateFriends(id, friends)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Print("friends updated with success")
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	dal, _ := dal.GetSession()

	users, err := dal.FindAll()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	payload := map[string]interface{}{
		"users": users,
	}
	log.Print("All users retrieved with success")
	json.NewEncoder(w).Encode(payload)
}

// our main function
func main() {

	dal = &Dal{}
	defer dal.CloseSession()

	router := mux.NewRouter()
	router.HandleFunc("/user", CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}/state", SaveState).Methods("PUT")
	router.HandleFunc("/user/{id}/state", GetState).Methods("GET")
	router.HandleFunc("/user/{id}/friends", AddFriends).Methods("PUT")
	router.HandleFunc("/user/{id}/friends", GetFriends).Methods("GET")
	router.HandleFunc("/user", GetAll).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
