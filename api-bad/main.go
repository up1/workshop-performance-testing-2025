// main.go
package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(db:3306)/demo")
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/users", getUsers)
	r.Run(":8080")
}

func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query error"})
		return
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		users = append(users, map[string]interface{}{
			"id": id, "name": name,
		})
	}

	c.JSON(http.StatusOK, users)
}
