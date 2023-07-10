package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type todos struct {
	ID   uint8  `gorm:"primaryKey;autoIncrement" form:"id"   json:"id"`
	Todo string `gorm:"not null"                 form:"todo" json:"todo"`
	Done bool   `gorm:"not null"                 form:"done" json:"done"`
}

func initDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func getTodosList(l *gin.Context) {
	db := initDb()
	var todos []todos
	db.Find(&todos)
	l.JSON(http.StatusOK, todos)
}

func addTodos(t *gin.Context) {
	db := initDb()

	var todo todos
	t.Bind(&todo)
	if todo.Todo != "" {
		db.Create(&todo)
		t.JSON(http.StatusOK, gin.H{"success": todo})
	} else {
		t.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Field is emty"})
	}
}

func getTodo(t *gin.Context) {
	db := initDb()
	id := t.Params.ByName("id")
	var todo todos
	db.First(&todo, id)
	if todo.ID != 0 {
		t.JSON(http.StatusOK, todo)
	} else {
		t.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	}
}

func updateTodo(t *gin.Context) {
	db := initDb()
	id := t.Params.ByName("id")
	var todo todos
	db.First(&todo, id)
	if todo.Todo != "" {
		if todo.ID != 0 {
			var newTodo todos
			t.Bind(&newTodo)
			result := todos{
				ID:   todo.ID,
				Todo: newTodo.Todo,
				Done: newTodo.Done,
			}
			db.Save(&result)
			t.JSON(http.StatusOK, gin.H{"success": result})
		} else {
			t.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		}
	} else {
		t.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Field is emty"})
	}
}

func deleteTodo(t *gin.Context) {
	db := initDb()
	id := t.Params.ByName("id")
	var todo todos
	db.First(&todo, id)
	if todo.ID != 0 {
		db.Delete(&todo)
		t.JSON(http.StatusOK, gin.H{"success": "Todo #" + id + "deleted"})
	} else {
		t.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
	}
}

func main() {
	router := gin.Default()
	db := initDb()
	db.AutoMigrate(&todos{})

	router.GET("/list", getTodosList)
	router.GET("/list/:id", getTodo)
	router.POST("/list/add", addTodos)
	router.DELETE("/list/delete/:id", deleteTodo)
	router.PUT("/list/update/:id", updateTodo)

	router.Run("localhost:8080")
}
