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
}

type WhiskeyCtrl struct {
	db *sqlx.DB
}

func (ctrl WhiskeyCtrl) All(c *gin.Context) {
	whiskeys := make([]Whiskey, 0)
	ctrl.db.Select(&whiskeys, `SELECT * FROM whiskeys ORDER BY distillery, name, age, size`)

	c.JSON(http.StatusOK, whiskeys)
}
