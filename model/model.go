package model

import (
	"crypto/md5"
	"database/sql"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	_ "github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/lib/pq"
)

func init() {
	gob.Register(&User{})
	gob.Register(&Whiskey{})
	gob.Register(&Post{})
}

type Model struct {
	Users    *Users
	Whiskeys *Whiskeys
	Posts    *Posts
}

func Connect(url string) (*Model, error) {
	conn, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return &Model{
		Users:    &Users{conn},
		Whiskeys: &Whiskeys{conn},
		Posts:    &Posts{conn},
	}, nil
}

type User struct {
	ID       string `db:"id"`
	EMail    string `db:"email"`
	Gravatar string `db:"gravatar"`
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

type Users struct {
	conn *sqlx.DB
}

type Whiskey struct {
	ID         string  `db:"id"`
	Distillery string  `db:"distillery"`
	Type       string  `db:"type"`
	Age        int     `db:"age"`
	ABV        float64 `db:"abv"`
	Size       float64 `db:"size"`
}

type Whiskeys struct {
	conn *sqlx.DB
}

func (srv *Whiskeys) All() ([]*Whiskey, error) {
	whiskeys := make([]*Whiskey, 0)
	err := srv.conn.Select(&whiskeys, "SELECT * FROM whiskeys")
	return whiskeys, err
}

type Post struct {
	ID      string    `db:"id"`
	User    string    `db:"user_id"`
	Whiskey string    `db:"whiskey_id"`
	Date    time.Time `db:"date"`
	Body    string    `db:"body"`
}

type Posts struct {
	conn *sqlx.DB
}

func (srv *Posts) FindByUser(user string, limit int) ([]*Post, error) {
	posts := make([]*Post, 0)
	err := srv.conn.Select(&posts, "SELECT * FROM posts WHERE user_id = $1 LIMIT $2", user, limit)
	return posts, err
}
