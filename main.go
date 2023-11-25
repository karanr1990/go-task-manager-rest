package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "postgres"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db

}

// add Task  represent task with its properties
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	Status      string    `json:"status"`
}

func getTasks(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	db := OpenConnection()
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var taskList []Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
		if err != nil {
			log.Fatal(err)
		}
		taskList = append(taskList, task)
	}

	ctx.IndentedJSON(http.StatusOK, taskList)
}
func addTask(ctx *gin.Context) {
	db := OpenConnection()
	var newTask Task

	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	stmt, err := db.Prepare("INSERT INTO tasks (id,title,description,duedate,status) VALUES ($1,$2,$3,$4,$5)")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(newTask.ID, newTask.Title, newTask.Description, newTask.DueDate, newTask.Status); err != nil {
		log.Fatal(err)
	}
	ctx.JSON(http.StatusOK, newTask)
}
func getTask(ctx *gin.Context) {
	db := OpenConnection()
	id := ctx.Param("id")
	var task Task
	results, err := db.Query("SELECT * FROM tasks WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}

	if results.Next() {
		err = results.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
		if err != nil {
			log.Fatal(err)
		}
	}
	ctx.JSON(http.StatusOK, task)

}
func updateTask(ctx *gin.Context) {
	db := OpenConnection()
	id := ctx.Param("id")

	var task Task
	row, err := db.Query("UPDATE tasks SET title=$2 ,description=$3 WHERE id=$1", id, &task.Title, &task.Description)
	if err != nil {
		log.Fatal(err)
	}
	if row.Next() {
		err := row.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Status)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx.JSON(http.StatusOK, row)
}

//func removeTask(ctx *gin.Context) {
//	id := ctx.Param("id")
//
//	for i, val := range tasks {
//		if val.ID == id {
//			tasks = append(tasks[:i], tasks[i+1:]...)
//			ctx.JSON(http.StatusOK, gin.H{"message": "Task removed"})
//			return
//		}
//	}
//
//	ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
//}

func main() {
	//gin is used to set up router engine for the API
	router := gin.Default()

	//get all tasks
	router.GET("/tasks", getTasks)

	//create new task
	router.POST("/tasks", addTask)

	//get specific task
	router.GET("/tasks/:id", getTask)

	//update task
	router.PUT("/tasks/:id", updateTask)

	//delete task
	//router.DELETE("/tasks/:id", removeTask)

	//listen and serve
	router.Run()

}
