package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Database(connString string) gin.HandlerFunc {
	db, err := sql.Open("postgres", os.Getenv("CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}
	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

func main() {
	godotenv.Load()

	router := gin.Default()

	router.Use(Database(os.Getenv("CONNECTION_STRING")))

	router.Run(":8080")
}
