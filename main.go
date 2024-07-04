package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	goutil "github.com/muhammadrivaldy/go-util"
	"github.com/redis/go-redis/v9"
)

var sqlDB *sql.DB
var rs *redsync.Redsync

func main() {

	var err error
	sqlDB, err = goutil.NewMySQL("root", "root", os.Getenv("DB_URL"), "automate-test", nil)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	pool := goredis.NewPool(redisClient)
	rs = redsync.New(pool)

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

	engine.Run(":8080")
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

	if err := rs.NewMutex("test-just"); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	var user User
	if err := c.Bind(&user); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	const sqlInsertUser = `INSERT INTO mst_users (name, address) VALUES (?, ?)`

	if _, err := sqlDB.Exec(sqlInsertUser, user.Name, user.Address); err != nil {
		fmt.Println(err)
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
