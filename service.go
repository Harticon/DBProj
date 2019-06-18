package DBproj

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
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

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err = s.access.CreateUser(usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (s *Service) SignIn(ctx echo.Context) error {

	//return login token

	var usr User

	err := ctx.Bind(&usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	query, err := s.access.GetUser(usr)
	if err != nil {
		//fmt.Println("error: ", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
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

	return ctx.JSON(http.StatusAccepted, Token{
		Token: t,
	})
}

func (s *Service) SetTask(ctx echo.Context) error {

	var tsk Task

	err := ctx.Bind(&tsk)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	var ok bool
	tsk.UserId, ok = ctx.Get("id").(int)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, err)
	}

	err = s.access.CreateTask(tsk)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, err)
}

func (s *Service) GetTaskByUserId(ctx echo.Context) error {

	from := ctx.QueryParam("from")
	to := ctx.QueryParam("to")

	f, ok := strconv.Atoi(from)
	if ok != nil {
		return ctx.JSON(http.StatusBadRequest, ok)
	}

	t, ok := strconv.Atoi(to)
	if ok != nil {
		return ctx.JSON(http.StatusBadRequest, ok)
	}

	result, err := s.access.GetTask(ctx.Get("id").(int), f, t)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, result)
}

func hashPassword(raw string) (string, error) {
	dk, err := scrypt.Key([]byte(raw), []byte("salt&pepper"), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dk), nil
}
