package main

import (
	"net/http"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type User struct {
	ID      string `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Nick    string `db:"nick" json:"nick"`
	EMail   string `db:"email" json:"email"`
	Picture string `db:"picture" json:"picture"`
}

type PrincipalCtrl struct {
	db *sqlx.DB
}

func (ctrl PrincipalCtrl) Get(c *gin.Context) {
	id, _ := c.Get("user")
	user := &User{}

	ctrl.db.Get(&user, `select * from users where id = $1`, id)
	c.JSON(http.StatusOK, user)
}

func (ctrl PrincipalCtrl) Create(c *gin.Context) {
	principal, _ := c.Get("user")

	user := &User{}
	c.BindJSON(user)
	user.ID = principal.(string)

	exists := false
	ctrl.db.Get(&exists, `select count(*) > 0 from users where id = $1`, principal)

	if exists {
		stmt, _ := ctrl.db.PrepareNamed(`update users set
		name = :name, nick = :nick, email = :email, picture = :picture
		where id = :id
		returning *`)
		stmt.Get(&user, user)
		c.JSON(http.StatusOK, user)

	} else {
		stmt, _ := ctrl.db.PrepareNamed(`insert into users 
		(id, name, nick, email, picture)
		values (:id, :name, :nick, :email, :picture)
		returning *`)
		stmt.Get(&user, user)
		c.JSON(http.StatusCreated, user)
	}
}
