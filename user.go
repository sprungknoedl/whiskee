package main

import (
	"net/http"
	"time"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type User struct {
	ID      string    `db:"id" json:"id"`
	Name    string    `db:"name" json:"name"`
	Nick    string    `db:"nick" json:"nick"`
	EMail   string    `db:"email" json:"email"`
	Picture string    `db:"picture" json:"picture"`
	Created time.Time `db:"created" json:"created"`
}

type PrincipalCtrl struct {
	db *sqlx.DB
}

func (this PrincipalCtrl) Get(c *gin.Context) {
	id, _ := c.Get("user")
	user := &User{}

	this.db.Get(&user, `select * from users where id = $1`, id)
	c.JSON(http.StatusOK, user)
}

func (this PrincipalCtrl) Create(c *gin.Context) {
	principal, _ := c.Get("user")

	user := &User{}
	c.BindJSON(user)
	user.ID = principal.(string)

	exists := false
	this.db.Get(&exists, `select count(*) > 0 from users where id = $1`, principal)

	if exists {
		stmt, _ := this.db.PrepareNamed(`update users set
		name = :name, nick = :nick, email = :email, picture = :picture, created = :created
		where id = :id
		returning *`)
		stmt.Get(&user, user)
		c.JSON(http.StatusOK, user)

	} else {
		stmt, _ := this.db.PrepareNamed(`insert into users 
		(id, name, nick, email, picture, created)
		values (:id, :name, :nick, :email, :picture, now())
		returning *`)
		stmt.Get(&user, user)
		c.JSON(http.StatusCreated, user)
	}
}

type UserCtrl struct {
	db *sqlx.DB
}

func (this UserCtrl) All(c *gin.Context) {
	rows := make([]User, 0)
	this.db.Select(&rows, `select * from users`)
	c.JSON(http.StatusOK, rows)
}

type FriendsCtrl struct {
	db *sqlx.DB
}

func (this FriendsCtrl) All(c *gin.Context) {
	principal, _ := c.Get("user")

	rows := make([]User, 0)
	this.db.Select(&rows, `select users.* from users
	join friends on (users.id = friends.a)
	where friends.b = $1`, principal)
	c.JSON(http.StatusOK, rows)
}

func (this FriendsCtrl) Create(c *gin.Context) {
	principal, _ := c.Get("user")
	friend := c.Param("id")

	this.db.Exec(`insert into friends (a, b) values ($1, $2)`, principal, friend)
	c.JSON(http.StatusCreated, friend)
}
