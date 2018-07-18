package main

import (
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

var settings = postgresql.ConnectionURL{
	Host:     "postgres",
	Database: "game",
	User:     "postgres",
	Password: "",
}

type Dal struct {
	Session sqlbuilder.Database
}

func (dal *Dal) GetSession() (*Dal, error) {
	// guarantees it'll only open one session
	if dal.Session != nil {
		return dal, nil
	}
	sess, err := postgresql.Open(settings)
	if err != nil {
		return nil, err
	}
	dal.Session = sess
	return dal, nil
}

func (dal *Dal) CloseSession() {
	if dal.Session != nil {
		dal.Session.Close()
	}
}

func (dal *Dal) InsertPlayer(newPlayer Player) error {
	players := dal.Session.Collection("player")
	_, err := players.Insert(newPlayer)
	return err
}

func (dal *Dal) UpdateState(score int, gamesPlayed int, id string) error {
	q := dal.Session.Update("player").Set(
		"score", score,
		"games", gamesPlayed,
	).Where("id = ?", id)
	_, err := q.Exec()
	return err
}

func (dal *Dal) GetState(player *Player, id string) error {
	players := dal.Session.Collection("player")
	res := players.Find("id", id)
	err := res.One(&player)
	return err
}

func (dal *Dal) GetFriends(id string) ([]Friend, error) {
	res, err := dal.Session.Query(`SELECT * FROM player WHERE ID IN (
			SELECT player1 FROM friendship WHERE player2 = ?
			UNION
			SELECT player2 FROM friendship WHERE player1 = ?
		)`, id, id)
	if err != nil {
		return nil, err
	}
	var friends []Friend
	iter := sqlbuilder.NewIterator(res)
	iter.All(&friends)

	return friends, nil
}

func (dal *Dal) UpdateFriends(id string, friends FriendsRequest) error {
	_, err := dal.Session.Exec(`DELETE FROM friendship WHERE player1 = ? OR player2 = ?`, id, id)
	if err != nil {
		return err
	}
	for _, friend := range *friends.IDs {
		// The ORM being used doesnt support bulk operations.
		// This operation will be expensive
		q := dal.Session.InsertInto("friendship").Columns("player1", "player2").Values(id, friend)
		_, err = q.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (dal *Dal) FindAll() ([]Player, error) {
	players := dal.Session.Collection("player")
	res := players.Find()
	var users []Player
	err := res.All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
