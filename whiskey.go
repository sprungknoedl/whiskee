package main

import (
	"net/http"

	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/sprungknoedl/whiskee/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

type Whiskey struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Distillery string  `json:"distillery"`
	Type       string  `json:"type"`
	Age        int     `json:"age"`
	ABV        float64 `json:"abv"`
	Size       float64 `json:"size"`
	Picture    string  `json:"picture"`
	Thumb      string  `json:"thumb"`
}

type WhiskeyCtrl struct {
	db *sqlx.DB
}

func (ctrl WhiskeyCtrl) All(c *gin.Context) {
	whiskeys := make([]Whiskey, 0)
	ctrl.db.Select(&whiskeys, `SELECT * FROM whiskeys ORDER BY distillery, name, age, size`)

	c.JSON(http.StatusOK, whiskeys)
}

func (this WhiskeyCtrl) Create(c *gin.Context) {
	whiskey := new(Whiskey)
	c.BindJSON(whiskey)

	this.db.Get(&whiskey.ID, `insert into whiskeys
	(distillery, name, type, age, abv, size, picture, thumb)
	values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`,
		whiskey.Distillery, whiskey.Name, whiskey.Type, whiskey.Age, whiskey.ABV,
		whiskey.Size, whiskey.Picture, whiskey.Thumb)

	c.JSON(http.StatusCreated, whiskey)
}
