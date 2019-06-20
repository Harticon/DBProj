package main

import (
	"fmt"
	"github.com/Harticon/DBproj"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
)

func main() {

	// todo CLI add something that's gonna print

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang, l",
			Value: "english",
			Usage: "Language for the greeting",
		},
		cli.StringFlag{
			Name:  "paymentMethod, p",
			Value: "cash",
			Usage: "Payment method",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) error {
				fmt.Println("Spoustime funkci complete")
				return nil
			},
		},
		{
			Name:    "begin",
			Aliases: []string{"b"},
			Usage:   "begin task",
			Action: func(c *cli.Context) error {
				fmt.Println("Spoustime funkci begin")
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Action = func(c *cli.Context) error {

		name := "moje nejlepší aplikace na světě"

		if c.NArg() > 0 {
			name = c.Args().Get(0)
			fmt.Printf("---------------------------%s\n",name)
		} else if c.String("lang") == "spanish" {
			fmt.Println("Hola", name)
		} else if c.String("lang") == "czech" {
			fmt.Println("Zdravíčko", name)
		} else {
			fmt.Println("Hello", name)
		}
		return nil


	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	govalidator.SetFieldsRequiredByDefault(true)

	//defer profile.Start().Stop()

	viper.SetDefault("db.conn", "prod.db")
	viper.SetDefault("secret", "secret")
	viper.SetDefault("hashSecret", "salt&peper")

	db, err := gorm.Open("sqlite3", viper.GetString("db.conn"))
	if err != nil {
		panic("failed to connect to database	")
	}

	db.AutoMigrate(&DBproj.User{}, &DBproj.Task{})

	access := DBproj.NewAccess(db)
	service := DBproj.NewService(access)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ug := e.Group("/auth")
	ug.POST("/signup", service.SignUp)
	ug.POST("/login", service.SignIn)

	tg := e.Group("/task")
	tg.Use(DBproj.UserMiddleware)
	tg.POST("/create", service.SetTask)
	tg.GET("/get", service.GetTaskByUserId)

	e.Logger.Fatal(e.Start(":8080"))

}
