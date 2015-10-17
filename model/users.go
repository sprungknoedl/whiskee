package model

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type User struct {
	ID       string `db:"id"`
	EMail    string `db:"email"`
	Gravatar string `db:"gravatar"`
}

type Users struct {
	conn *sqlx.DB
}

func (srv *Users) All() ([]*User, error) {
	users := make([]*User, 0)
	err := srv.conn.Select(&users, `SELECT * FROM users`)
	return users, err
}

func (srv *Users) Friends(id string) ([]*User, error) {
	users := make([]*User, 0)
	err := srv.conn.Select(&users,
		`SELECT u.* FROM users u
		JOIN friends f ON u.id = f.b
		WHERE f.a = $1`, id)
	return users, err
}

func (srv *Users) Get(id string) (*User, error) {
	u := &User{}
	err := srv.conn.Get(u, "SELECT * FROM users WHERE id = $1", id)
	return u, err
}

func (srv *Users) Create(id, email string) (*User, error) {
	data := []byte(strings.TrimSpace(email))
	gravatar := fmt.Sprintf("%x", md5.Sum(data))

	user := &User{
		ID:       id,
		EMail:    email,
		Gravatar: gravatar,
	}

	_, err := srv.conn.NamedExec("INSERT INTO users (id, email, gravatar) VALUES (:id, :email, :gravatar);", user)
	return user, err
}

func (srv *Users) GetOrCreate(id, email string) (*User, error) {
	user, err := srv.Get(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// register user
	if err == sql.ErrNoRows {
		user, err = srv.Create(id, email)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (srv *Users) AddFriend(a, b string) error {
	_, err := srv.conn.Exec(`INSERT into friends (a, b) VALUES ($1, $2)`, a, b)
	return err
}
