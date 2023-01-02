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

type student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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

	router.GET("/students", getStudents)
	router.POST("/students", createStudent)
	router.PATCH("/students/:id", updateStudent)
	router.DELETE("/students/:id", deleteStudent)

	router.Run(":8080")
}

func getStudents(c *gin.Context) {
	db := c.MustGet("DB").(*sql.DB)
	rows, err := db.Query("SELECT id,name FROM students")
	var students []student
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Students not found"})
	}

	for rows.Next() {
		var name string
		var id int

		err = rows.Scan(&id, &name)

		if err != nil {
			log.Fatal(err)
		}

		student := student{ID: id, Name: name}
		students = append(students, student)
	}
	c.JSON(http.StatusOK, students)
}

func createStudent(c *gin.Context) {
	db := c.MustGet("DB").(*sql.DB)
	insertStmt := `INSERT INTO "students"("name") VALUES($1)`

	var newStudent student

	if err := c.BindJSON(&newStudent); err != nil {
		log.Fatal(err)
		return
	}

	_, err := db.Exec(insertStmt, newStudent.Name)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created new student"})
}

func updateStudent(c *gin.Context) {
	db := c.MustGet("DB").(*sql.DB)
	id := c.Param("id")
	updateStmt := `UPDATE "students" SET "name"=$1 WHERE "id"=$2`
	var updatedStudent student

	if err := c.BindJSON(&updatedStudent); err != nil {
		log.Fatal(err)
		return
	}

	_, err := db.Exec(updateStmt, updatedStudent.Name, id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated student"})
}

func deleteStudent(c *gin.Context) {
	db := c.MustGet("DB").(*sql.DB)
	id := c.Param("id")
	deleteStmt := `DELETE FROM "students" WHERE "id"=$1`

	_, err := db.Exec(deleteStmt, id)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted student"})
}
