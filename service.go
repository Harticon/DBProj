package DBproj

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"net/http"
	"regexp"
)

type Service struct {
	access IAccesser
}

func NewService(access IAccesser) *Service {
	return &Service{
		access: access,
	}
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

func (s *Service) SignUp(ctx echo.Context) error {

	var usr User

	err := ctx.Bind(&usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if usr.Email == "" || usr.Password == "" {
		return ctx.JSON(http.StatusBadRequest, errors.New("empty password and email"))
	}

	//Validate email
	if !validateEmail(usr.Email) {
		fmt.Println("email not valid")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err = s.access.CreateUser(usr)
	return ctx.JSON(http.StatusCreated, err)
}

func (s *Service) SignIn(ctx echo.Context) error {

	//return login token

	var usr User

	err := ctx.Bind(&usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	query := s.access.GetUser(usr)
	if query == (User{}) {
		return ctx.JSON(http.StatusBadRequest, errors.New("user does not exists"))
	}

	//create token

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = query.Email
	claims["id"] = query.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusAccepted, map[string]string{
		"token": t,
	})
}

func (s *Service) SetTask(ctx echo.Context) error {

	var tsk Task

	err := ctx.Bind(&tsk)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	tsk.UserId = ctx.Get("id").(int)

	s.access.CreateTask(tsk)
	return ctx.JSON(http.StatusCreated, err)
}

func (s *Service) GetTaskByUserId(ctx echo.Context) error {

	from := ctx.QueryParam("from")
	to := ctx.QueryParam("to")

	f, ok := strconv.Atoi(from)
	if ok != nil {
		fmt.Printf("Nebylo zadano cislo")
		return ctx.JSON(http.StatusBadRequest, ok)
	}

	t, ok := strconv.Atoi(to)
	if ok != nil {
		fmt.Printf("Nebylo zadano cislo")
		return ctx.JSON(http.StatusBadRequest, ok)
	}

	result := s.access.GetTask(ctx.Get("id").(int), f, t)

	return ctx.JSON(http.StatusOK, result)
}
