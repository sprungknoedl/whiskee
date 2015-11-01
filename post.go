package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type Post struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	Body     string    `json:"body"`
	Security string    `json:"security"`
	User     *User     `json:"user"`
	Whiskey  *Whiskey  `json:"whiskey"`
}

type postEntity struct {
	PostID            int       `db:"post_id"`
	PostBody          string    `db:"post_body"`
	PostDate          time.Time `db:"post_date"`
	PostSecurity      string    `db:"post_security"`
	UserID            string    `db:"user_id"`
	UserName          string    `db:"user_name"`
	UserNick          string    `db:"user_nick"`
	UserEMail         string    `db:"user_email"`
	UserPicture       string    `db:"user_picture"`
	WhiskeyID         int       `db:"whiskey_id"`
	WhiskeyDistillery string    `db:"whiskey_distillery"`
	WhiskeyName       string    `db:"whiskey_name"`
	WhiskeyType       string    `db:"whiskey_type"`
	WhiskeyAge        int       `db:"whiskey_age"`
	WhiskeyABV        float64   `db:"whiskey_abv"`
	WhiskeySize       float64   `db:"whiskey_size"`
}

func (p postEntity) Post() Post {
	user := &User{
		ID:      p.UserID,
		Name:    p.UserName,
		Nick:    p.UserNick,
		EMail:   p.UserEMail,
		Picture: p.UserPicture,
	}

	whiskey := &Whiskey{
		ID:         p.WhiskeyID,
		Distillery: p.WhiskeyDistillery,
		Name:       p.WhiskeyName,
		Type:       p.WhiskeyType,
		Age:        p.WhiskeyAge,
		ABV:        p.WhiskeyABV,
		Size:       p.WhiskeySize,
	}

	return Post{
		ID:       p.PostID,
		Date:     p.PostDate,
		Body:     p.PostBody,
		Security: p.PostSecurity,
		User:     user,
		Whiskey:  whiskey,
	}
}

type PostCtrl struct {
	db *sqlx.DB
}

// GET    /:type        -> All()
// GET    /:type/:id    -> One(id)
// POST   /:type        -> Create(entity)
// PUT    /:type/:id    -> Update(id, entity)
// DELETE /:type/:id    -> Delete(id)

func (ctrl PostCtrl) All(c *gin.Context) {
	principal, _ := c.Get("user")

	rows := make([]postEntity, 0)
	ctrl.db.Select(&rows, `select 
	p.id as post_id, p.body as post_body, p.date as post_date, p.security as post_security,
	u.id as user_id, u.name as user_name, u.nick as user_nick, u.email as user_email, u.picture as user_picture,
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
	limit $2`, principal, 20)

	posts := make([]Post, len(rows))
	for i, row := range rows {
		posts[i] = row.Post()
	}

	c.JSON(http.StatusOK, posts)
}

func (ctrl PostCtrl) Create(c *gin.Context) {
	principal, _ := c.Get("user")

	post := new(Post)
	c.BindJSON(post)

	var id int
	ctrl.db.Get(&id, `insert into posts
	(body, date, security, user_id, whiskey_id)
	values ($1, $2, $3, $4, $5) returning id`, post.Body, post.Date, post.Security, principal, post.Whiskey.ID)

	var entity postEntity
	ctrl.db.Get(&entity, `select 
	p.id as post_id, p.body as post_body, p.date as post_date, p.security as post_security,
	u.id as user_id, u.name as user_name, u.nick as user_nick, u.email as user_email, u.picture as user_picture,
	w.id as whiskey_id, w.distillery as whiskey_distillery, w.name as whiskey_name,
		w.type as whiskey_type, w.age as whiskey_age, w.abv as whiskey_abv, 
		w.size as whiskey_size

	from posts p
	join users u on u.id = p.user_id
	join whiskeys w on w.id = p.whiskey_id

	where p.id = $1`, id)
	c.JSON(http.StatusCreated, entity.Post())
}

func (ctrl PostCtrl) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	principal, _ := c.Get("user")

	ctrl.db.Exec(`delete from posts where id = $1 and user_id = $2`, id, principal)
	c.Data(http.StatusNoContent, gin.MIMEJSON, nil)
}
