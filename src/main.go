package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

var settings = postgresql.ConnectionURL{
	Host:     "172.17.0.2",
	Database: "game",
	User:     "postgres",
	Password: "",
}

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
	sess, _ := postgresql.Open(settings)
	defer sess.Close()
	players := sess.Collection("player")

	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	newPlayer := Player{
		ID:   uuid.New(),
		Name: player.Name,
	}

	_, _ = players.Insert(newPlayer)

	json.NewEncoder(w).Encode(newPlayer)
}

func SaveState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	sess, _ := postgresql.Open(settings)

	defer sess.Close()

	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	q := sess.Update("player").Set(
		"score", player.Score,
		"games", player.GamesPlayed,
	).Where("id = ?", id)

	_, err = q.Exec()
}

func GetState(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]
	sess, _ := postgresql.Open(settings)

	defer sess.Close()
	players := sess.Collection("player")
	var player Player

	res := players.Find("id", id)
	err := res.One(&player)
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
	sess, _ := postgresql.Open(settings)

	defer sess.Close()

	res, err := sess.Query(`SELECT * FROM player WHERE ID IN (
			SELECT player1 FROM friendship WHERE player2 = ?
			UNION
			SELECT player2 FROM friendship WHERE player1 = ?
		)`, id, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var friends []Friend
	iter := sqlbuilder.NewIterator(res)
	iter.All(&friends)

	payload := map[string]interface{}{
		"friends": friends,
	}
	json.NewEncoder(w).Encode(payload)
}

func AddFriends(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	sess, _ := postgresql.Open(settings)

	defer sess.Close()

	var friends FriendsRequest
	err := json.NewDecoder(r.Body).Decode(&friends)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = sess.Exec(`DELETE FROM friendship WHERE player1 = ? OR player2 = ?`, id, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for _, friend := range *friends.IDs {
		// No bulk operations in this lib. Found it too late
		q := sess.InsertInto("friendship").Columns("player1", "player2").Values(id, friend)
		_, err = q.Exec()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

}

func GetAll(w http.ResponseWriter, r *http.Request) {
	var users []Player
	sess, _ := postgresql.Open(settings)
	defer sess.Close()

	players := sess.Collection("player")
	res := players.Find()
	err := res.All(&users)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	payload := map[string]interface{}{
		"users": users,
	}
	json.NewEncoder(w).Encode(payload)
}

// our main function
func main() {

	router := mux.NewRouter()
	router.HandleFunc("/user", CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}/state", SaveState).Methods("PUT")
	router.HandleFunc("/user/{id}/state", GetState).Methods("GET")
	router.HandleFunc("/user/{id}/friends", AddFriends).Methods("PUT")
	router.HandleFunc("/user/{id}/friends", GetFriends).Methods("GET")
	router.HandleFunc("/user", GetAll).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
