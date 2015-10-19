package model

import (
	"encoding/gob"
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

type Whiskey struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
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
	err := srv.conn.Select(&whiskeys, `SELECT * FROM whiskeys ORDER BY distillery, name, age, size`)
	return whiskeys, err
}

func (srv *Whiskeys) Create(whiskey *Whiskey) (*Whiskey, error) {
	_, err := srv.conn.NamedExec(`
		INSERT INTO whiskeys (distillery, name, type, age, abv, size) 
		VALUES (:distillery, :name, :type, :age, :abv, :size);`, whiskey)
	return whiskey, err
}

type Post struct {
	ID       int
	User     *User
	Whiskey  *Whiskey
	Date     time.Time
	Body     string
	Security string
}

type Posts struct {
	conn *sqlx.DB
}

func (srv *Posts) GetNewsFeed(user string, limit int) ([]*Post, error) {
	rows := make([]struct {
		PostID            int       `db:"post_id"`
		PostBody          string    `db:"post_body"`
		PostDate          time.Time `db:"post_date"`
		PostSecurity      string    `db:"post_security"`
		UserID            string    `db:"user_id"`
		UserEMail         string    `db:"user_email"`
		UserGravatar      string    `db:"user_gravatar"`
		WhiskeyID         int       `db:"whiskey_id"`
		WhiskeyDistillery string    `db:"whiskey_distillery"`
		WhiskeyName       string    `db:"whiskey_name"`
		WhiskeyType       string    `db:"whiskey_type"`
		WhiskeyAge        int       `db:"whiskey_age"`
		WhiskeyABV        float64   `db:"whiskey_abv"`
		WhiskeySize       float64   `db:"whiskey_size"`
	}, 0)

	err := srv.conn.Select(&rows, `select 
	p.id as post_id, p.body as post_body, p.date as post_date, p.security as post_security,
	u.id as user_id, u.email as user_email, u.gravatar as user_gravatar,
	w.id as whiskey_id, w.distillery as whiskey_distillery, w.name as whiskey_name,
		w.type as whiskey_type, w.age as whiskey_age, w.abv as whiskey_abv, 
		w.size as whiskey_size

	from posts p
	join users u on u.id = p.user_id
	join whiskeys w on w.id = p.whiskey_id

	where
		-- select my posts
		p.id in (select p.id from posts p
			where p.user_id = $1) or

		-- select public posts of my friends
		p.id in (select p.id from posts p 
			join friends f on p.user_id = f.b
			where (f.a = $1 and p.security = 'public')) or

		-- select private posts of my friends
		p.id in (with mutual as (
				select b as id from friends where a = $1 intersect 
				select a as id from friends where b = $1)
			select p.id from posts p
			join mutual m on m.id = p.user_id
			where p.security = 'friends')

	order by p.date desc
	limit $2`, user, limit)

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
			Name:       row.WhiskeyName,
			Type:       row.WhiskeyType,
			Age:        row.WhiskeyAge,
			ABV:        row.WhiskeyABV,
			Size:       row.WhiskeySize,
		}

		posts[i] = &Post{
			ID:       row.PostID,
			Date:     row.PostDate,
			Body:     row.PostBody,
			Security: row.PostSecurity,
			User:     user,
			Whiskey:  whiskey,
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
