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

func (srv *Users) All() ([]*User, error) {
	users := make([]*User, 0)
	err := srv.conn.Select(&users, `SELECT * FROM users`)
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

type Users struct {
	conn *sqlx.DB
}

type Whiskey struct {
	ID         int     `db:"id"`
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
	err := srv.conn.Select(&whiskeys, `SELECT * FROM whiskeys`)
	return whiskeys, err
}

func (srv *Whiskeys) Create(whiskey *Whiskey) (*Whiskey, error) {
	_, err := srv.conn.NamedExec(`
		INSERT INTO whiskeys (distillery, type, age, abv, size) 
		VALUES (:distillery, :type, :age, :abv, :size);`, whiskey)
	return whiskey, err
}

type Post struct {
	ID      int
	User    *User
	Whiskey *Whiskey
	Date    time.Time
	Body    string
}

type Posts struct {
	conn *sqlx.DB
}

func (srv *Posts) GetNewsFeed(user string, limit int) ([]*Post, error) {
	rows := make([]struct {
		PostID            int       `db:"post_id"`
		PostBody          string    `db:"post_body"`
		PostDate          time.Time `db:"post_date"`
		UserID            string    `db:"user_id"`
		UserEMail         string    `db:"user_email"`
		UserGravatar      string    `db:"user_gravatar"`
		WhiskeyID         int       `db:"whiskey_id"`
		WhiskeyDistillery string    `db:"whiskey_distillery"`
		WhiskeyType       string    `db:"whiskey_type"`
		WhiskeyAge        int       `db:"whiskey_age"`
		WhiskeyABV        float64   `db:"whiskey_abv"`
		WhiskeySize       float64   `db:"whiskey_size"`
	}, 0)

	err := srv.conn.Select(&rows, `
	WITH undirected (a, b) AS (
		SELECT a, b FROM friends
		UNION ALL
		SELECT b, a FROM friends)
	SELECT 
		p.id as post_id, p.body as post_body, p.date as post_date,
		u.id as user_id, u.email as user_email, u.gravatar as user_gravatar,
		w.id as whiskey_id, w.distillery as whiskey_distillery, w.type as whiskey_type,
			w.age as whiskey_age, w.abv as whiskey_abv, w.size as whiskey_size
		FROM posts p
		JOIN whiskeys w ON (p.whiskey_id = w.id)
		JOIN users u ON (p.user_id = u.id)
		WHERE p.user_id = $1 OR 
			p.user_id IN (SELECT b FROM undirected WHERE a = $1)
		ORDER BY p.date DESC LIMIT $2`, user, limit)

	posts := make([]*Post, len(rows))
	for i, row := range rows {
		user := &User{
			ID:       row.UserID,
			EMail:    row.UserEMail,
			Gravatar: row.UserGravatar,
		}

		whiskey := &Whiskey{
			ID:         row.WhiskeyID,
			Distillery: row.WhiskeyDistillery,
			Type:       row.WhiskeyType,
			Age:        row.WhiskeyAge,
			ABV:        row.WhiskeyABV,
			Size:       row.WhiskeySize,
		}

		posts[i] = &Post{
			ID:      row.PostID,
			Date:    row.PostDate,
			Body:    row.PostBody,
			User:    user,
			Whiskey: whiskey,
		}
	}

	return posts, err
}

func (srv *Posts) Create(post *Post) (*Post, error) {
	_, err := srv.conn.Exec(`
		INSERT INTO posts (user_id, whiskey_id, date, body) 
		VALUES ($1, $2, $3, $4);`,
		post.User.ID, post.Whiskey.ID, post.Date, post.Body)
	return post, err
}

func (srv *Posts) Delete(id int) error {
	_, err := srv.conn.Exec(`DELETE from posts where id = $1`, id)
	return err
}
