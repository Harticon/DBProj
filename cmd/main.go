package main

import (
	"github.com/Harticon/DBproj"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func main() {

	viper.SetDefault("db.conn", "test.db")

	db, err := gorm.Open("sqlite3", viper.GetString("db.conn"))
	if err != nil {
		panic("failed to connect to database	")
	}

	db.AutoMigrate(&DBproj.User{})
	db.AutoMigrate(&DBproj.Task{})

	access := DBproj.NewAccess(db)
	service := DBproj.NewService(access)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ug := e.Group("/auth")
	ug.POST("/signup", service.SignUp)
	ug.POST("/login", service.SignIn)

	tg := e.Group("/task")
	tg.POST("/create", service.SetTask)
	tg.GET("/get", service.GetTaskByUserId)

	e.Logger.Fatal(e.Start(":8080"))

}
