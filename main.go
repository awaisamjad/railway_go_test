package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

// Initialize the SQLite database connection
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create the users table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

// Route handler to add a user to the database
func addUser(ctx *gin.Context) {
	// Hardcoded user data to be added
	name := "John Doe"

	// Insert the user into the database
	query := "INSERT INTO users (name) VALUES (?)"
	_, err := db.Exec(query, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}

	// Redirect back to the home route after adding the user
	ctx.Redirect(http.StatusFound, "/")
}

// Route handler to fetch and display all users
func getUsers(ctx *gin.Context) {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []gin.H

	// Iterate over rows and append them to the users slice
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		users = append(users, gin.H{"id": id, "name": name})
	}

	// Display the users as JSON
	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func main() {
	// Initialize the database
	initDB()

	// Create a new Gin router
	r := gin.Default()

	// Define the home route
	r.GET("/", getUsers)

	// Define the route for adding a user
	r.GET("/add", addUser)

	// Start the server
	r.Run(":8080")
}

