package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type student struct {
	Name string  `json:"name"`
	ID   string  `json:"id"`
	Age  int     `json:"age"`
	Gpa  float32 `json:"gpa"`
}

var students = []student{
	{Name: "John Doe", ID: "2", Age: 12, Gpa: 2.5},
	{Name: "Nicole Jame", ID: "7", Age: 15, Gpa: 1.5},
	{Name: "Sheila Smith", ID: "5", Age: 22, Gpa: 4.5},
	{Name: "Daniel Jones", ID: "4", Age: 16, Gpa: 2.8},
	{Name: "Smith Roe", ID: "1", Age: 10, Gpa: 3.5},
}

func addStudent(context *gin.Context) {
	var newStudent student
	if err := context.BindJSON(&newStudent); err != nil {
		return
	}
	students = append(students, newStudent)
	context.IndentedJSON(http.StatusCreated, newStudent)
}

func getStudents(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, students)
}

func getStudent(id string) (*student, error) {
	for i, b := range students {
		if b.ID == id {
			return &students[i], nil
		}

	}
	return nil, errors.New("student not found")
}

func getStudentByGrade(c *gin.Context) {
	id := c.Param("id")
	student, err := getStudent(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "student with such grade does not exist"})
		return
	}
	c.IndentedJSON(http.StatusOK, student)
}

func reduceGpa(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error in getting query params": "could not get params"})
		return
	}
	student, err := getStudent(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error in getting student": "could not found student"})
		return
	}
	if student.Gpa <= 1 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error in student gpa number": "gpa too low to be reduced"})
		return
	}
	student.Gpa -= 1
	c.IndentedJSON(http.StatusOK, student)
}

func main() {
	router := gin.Default()

	router.GET("/students", getStudents)
	router.POST("/students", addStudent)
	router.GET("/students/:grade", getStudentByGrade)
	router.PATCH("/reducegpa", reduceGpa)

	router.Run("localhost:8080")

}
