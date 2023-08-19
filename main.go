package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int    `gorm:"<-:false" json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Age      string `json:"age" binding:"required"`
}

var db *gorm.DB
var err error

func Insert(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	db.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "user created successfully", "data": user})
}

func SelectAll(c *gin.Context) {
	var user []User
	db.Find(&user)
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func Select(c *gin.Context) {
	var user User
	id := c.Param("id")
	db.First(&user, id)

	if user.ID != 0 {
		c.JSON(http.StatusCreated, gin.H{
			"message": "user found",
		})
		c.JSON(http.StatusCreated, gin.H{"data": user})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
	}
}

func Delete(c *gin.Context) {
	var user User
	id := c.Param("id")
	db.First(&user, id)

	if user.ID != 0 {
		db.Where("id = ?", id).Delete(&user)
		c.JSON(http.StatusCreated, gin.H{
			"message": "user deleted",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
	}
}

func Update(c *gin.Context) {
	var o_user User
	id := c.Param("id")
	db.First(&o_user, id)

	if o_user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	var n_user User
	if err := c.ShouldBindJSON(&n_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	o_user.ID = n_user.ID
	o_user.Name = n_user.Name
	o_user.Email = n_user.Email
	o_user.Password = n_user.Password
	o_user.Age = n_user.Age
	db.Save(&o_user)
	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"data":    o_user,
	})
}

func Login(c *gin.Context) {
	var user User
	email := c.PostForm("email")
	password := c.PostForm("password")

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not a registered user"})
		return
	}

	if user.Password == password {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome " + user.Name})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
	}
}

func main() {

	dsn := "host=localhost user=mamdouhhazem dbname=postgres sslmode = disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.AutoMigrate(&User{})

	r := gin.Default()
	r.POST("/InsertUser", Insert)
	r.DELETE("/DeleteUser/:id", Delete)
	r.GET("/SelectUser", SelectAll)
	r.GET("/SelectAllUsers/:id", Select)
	r.PATCH("/UpdateUser/:id", Update)
	r.POST("/Login", Login)
	r.Run(":8080")
}
