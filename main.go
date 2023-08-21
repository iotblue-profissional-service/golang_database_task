package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	UserID   int `gorm:"unique_index"`
	Name     string
	Email    string `gorm:"unique_index"`
	Password string
	Age      int
}

var db *gorm.DB
var err error

// ConnectDB the connection of database
func ConnectDB() {
	db, err = gorm.Open("postgres", "host = 'localhost' port = 5432  user = postgres dbname = FirstTask password = maryam sslmode=disable ")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("database is connected")
	}

}

// the entry of the EndPoint
func main() {
	// start your task
	router := gin.Default()
	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong pong",
		})
	})

	ConnectDB()
	db.AutoMigrate(&User{})

	//Routers
	router.POST("/createUser", createUser)
	router.PUT("/updateUser/:id", updateUser)
	router.GET("/getUser/:id", getUser)
	router.GET("/getAllUsers", getAllUsers)
	router.DELETE("/deleteUser/:id", deleteUser)
	router.POST("/HomePage", login)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

// 1- Create User
func createUser(ctx *gin.Context) {
	body := User{}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User is not defined")
		return

	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Request")

	}
	res := db.Create(&body)
	if res.Error != nil {
		ctx.AbortWithStatusJSON(400, "can't insert new user")
	} else {
		ctx.JSON(http.StatusOK, "user inserted successfully")

	}
}

// 2- Update User
func updateUser(ctx *gin.Context) {
	var user User
	var UpUser User
	userID := ctx.Param("id")
	if err := db.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "user is not found"})
		return
	}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "invalid data..User is not defined")
		return
	}
	err = json.Unmarshal(data, &UpUser)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Request")
		return
	}

	db.Model(&User{}).Updates(UpUser)

	ctx.JSON(http.StatusOK, "user updated successfully")
}

// 3- Get User
func getUser(ctx *gin.Context) {
	userm := ctx.Param("id")
	var user User
	res := db.First(&user, userm)
	if res.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)

}

// 4- Get all Users
func getAllUsers(ctx *gin.Context) {
	var user []User
	res := db.Find(&user)
	if res.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "users not found"})
		return
	}
	ctx.JSON(http.StatusOK, res)

}

// 5- Delete User
func deleteUser(ctx *gin.Context) {
	var user User
	userID := ctx.Param("id")
	if err := db.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": "user is not found"})
		return
	}
	db.Delete(&user)
	ctx.JSON(http.StatusOK, "user deleted successfully")
}

// router.POST("/HomePage", login)
// 6- Login User (using email & password returning a message welcome [username] )
func login(ctx *gin.Context) {
	var user User
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Bad Request"})
		return
	}

	res := db.Where("email =? AND password =?  ", loginData.Email, loginData.Password).First(&user)
	if res.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"Error": " Can't Login ... Try Again"})
		return
	}
	ctx.JSON(http.StatusOK, "user Logged successfully")
}
