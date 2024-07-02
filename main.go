package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	goutil "github.com/muhammadrivaldy/go-util"
)

var sqlDB *sql.DB

func main() {

	engine := gin.Default()

	var err error
	sqlDB, err = goutil.NewMySQL("root", "root", os.Getenv("DB_URL"), "automate-test", nil)
	if err != nil {
		panic(err)
	}

	engine.GET("/health/check", func(c *gin.Context) { c.JSON(http.StatusOK, nil) })
	engine.POST("/users", handlerCreateUser)
	engine.GET("/users", handlerGetUsers)

	engine.Run(":80")
}

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Users struct {
	Users []User `json:"users"`
}

func handlerCreateUser(c *gin.Context) {

	var user User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	const sqlInsertUser = `INSERT INTO mst_users (name, address) VALUES (?, ?)`

	if _, err := sqlDB.Exec(sqlInsertUser, user.Name, user.Address); err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func handlerGetUsers(c *gin.Context) {

	const sqlSelectUsers = `SELECT id, name, address FROM mst_users`

	rows, err := sqlDB.Query(sqlSelectUsers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	defer rows.Close()

	var users Users

	for rows.Next() {

		var user User

		if err := rows.Scan(&user.ID, &user.Name, &user.Address); err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		users.Users = append(users.Users, user)
	}

	c.JSON(http.StatusOK, users)
}
