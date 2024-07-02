package main

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	goutil "github.com/muhammadrivaldy/go-util"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

var sqlDB *sql.DB

func main() {

	var err error
	sqlDB, err = goutil.NewMySQL("root", "root", os.Getenv("DB_URL"), "automate-test", nil)
	if err != nil {
		panic(err)
	}

	processName := os.Args[1]

	if processName == "service" {
		runService()
	} else if processName == "migration" {
		runMigration()
	} else {
		panic(errors.New("your request does not valid, please enter the correct parameter (service or migration)"))
	}
}

func runMigration() {

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"automate-test",
		driver)
	if err != nil {
		panic(err)
	}

	m.Up()

}

func runService() {

	engine := gin.Default()

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
