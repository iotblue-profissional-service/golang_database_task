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
	CompanyID int
	Name      string `gorm:"type:varchar(150)"`
	Email     string `gorm:"type:varchar(150);unique_index"`
	Password  string
	Age       int
}

var db *gorm.DB
var err error

func ConnectDBGorm() {
	db, err = gorm.Open("postgres", "host='localhost' port=5432 user=postgres dbname=task1 password=200315 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("success a5eeran")
	}
}
func main() {
	route := gin.Default()
	route.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	ConnectDBGorm()
	db.AutoMigrate(&User{})
	route.POST("/createuser", CreateUser) // Add this line
	route.PUT("/updateuser", UpdateUser)
	route.GET("/getuser/:{id}", GetUser)
	route.GET("/getusers", GetUsers)
	route.DELETE("/deleteuser/:{id}", DeleteUser)
	route.GET("/", Login)
	err := route.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
	// Add this line

}

/*
	func GetUser(ctx *gin.Context) {
		var user User
		id := ctx.Param("id")
		res := db.First(&user, id)
		if res.Error != nil {
			if user.ID == 0 {
				ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found! Try again with correct ID"})
				return
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "you are doing something wrong"})
			}
		}
	}
*/
func GetUser(ctx *gin.Context) {
	var user User
	id := ctx.Param("id")
	result := db.First(&user, id)
	if result.Error != nil {
		if result.RecordNotFound() {
			ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Retrieval failed"})
		}
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": user})
}
func GetUsers(ctx *gin.Context) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		if result.RecordNotFound() {
			ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Retrieval failed"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": users})
}
func CreateUser(ctx *gin.Context) {
	body := User{}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "User is not defined")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	result := db.Create(&body)
	if result.Error != nil {
		fmt.Println(result.Error)
		ctx.AbortWithStatusJSON(400, "Couldn't create the new user.")
	} else {
		ctx.JSON(http.StatusOK, "User is successfully created.")
	}
}
func UpdateUser(ctx *gin.Context) {
	var user User
	var updatedUser User
	id := ctx.Param("id")
	db.First(&user, id)
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		return
	}
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	err = json.Unmarshal(data, &updatedUser)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	db.Model(&user).Updates(updatedUser)
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "User updated successfully!"})
}
func DeleteUser(ctx *gin.Context) {
	var user User
	id := ctx.Param("id")
	res := db.First(&user, id)
	if res.Error != nil {
		if res.RecordNotFound() {
			ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Retrieval failed"})
		}
		return
	}
	db.Delete(&user)
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "User deleted successfully!"})
}
func Login(ctx *gin.Context) {
	var user User
	var loginUser User
	data, err := ctx.GetRawData()
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	err = json.Unmarshal(data, &loginUser)
	if err != nil {
		ctx.AbortWithStatusJSON(400, "Bad Input")
		return
	}
	result := db.Where("email = ? AND password = ?", loginUser.Email, loginUser.Password).First(&user)
	//TALAMA EL EMAIL UNIQUE MOSTA7EEL FE 2 USERS Y5O4O MAKAN BA3D
	if result.Error != nil {
		if result.RecordNotFound() {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid email or password, GO TO UPDATE USER TO CHANGE PASSWORD"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": fmt.Sprintf("Welcome back !!!!! <3 also end of task :) %s", user.Name)})
}
