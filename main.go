package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Age      int
	ID       int
}

func main() {

	dsn := "host=localhost user=postgres password=Victor@1939 dbname=Task1 port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to connect database")
	}
	println("Please Choose TO create User press 1, To Update User press 2, To Delete User press 3, To Get User press 4")
	var choice int
	fmt.Scanln(&choice)
	if choice == 1 {
		println("Enter Name")
		var name string
		fmt.Scanln(&name)
		println("Enter Email")
		var email string
		fmt.Scanln(&email)
		println("Enter Password")
		var password string
		fmt.Scanln(&password)
		println("Enter Age")
		var age int
		fmt.Scanln(&age)
		db.Create(&User{Name: name, Email: email, Password: password, Age: age})
	}

	if choice == 2 {
		println("Enter Name")
		var name string
		fmt.Scanln(&name)
		println("Enter Email")
		var email string
		fmt.Scanln(&email)
		println("Enter Password")
		var password string
		fmt.Scanln(&password)
		println("Enter Age")
		var age int
		fmt.Scanln(&age)
		db.Model(&User{}).Where("name = ?", name).Updates(User{Name: name, Email: email, Password: password, Age: age})
	}

	if choice == 3 {
		println("Enter Name")
		var name string
		fmt.Scanln(&name)
		db.Where("name = ?", name).Delete(&User{})
	}
}
