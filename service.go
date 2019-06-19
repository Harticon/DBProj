package DBproj

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/scrypt"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"net/http"
)

type Service struct {
	access IAccesser
}

func NewService(access IAccesser) *Service {
	return &Service{
		access: access,
	}
}

func (s *Service) SignUp(ctx echo.Context) error {

	var usr User

	err := ctx.Bind(&usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if usr.Email == "" || usr.Password == "" {
		return ctx.JSON(http.StatusBadRequest, errors.New("empty password or email"))
	}

	_, err = govalidator.ValidateStruct(usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	RetUsr, err := s.access.CreateUser(usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, RetUsr)
}

func (s *Service) SignIn(ctx echo.Context) error {

	//return login token

	var usr User

	err := ctx.Bind(&usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	_, err = govalidator.ValidateStruct(usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	usr.Password, err = hashPassword(usr.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	query, err := s.access.GetUser(usr)
	if err != nil {
		fmt.Println("error: ", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = query.Email
	claims["id"] = query.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(viper.GetString("secret")))
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
		return ctx.JSON(http.StatusUnauthorized, errors.New("couldnt resolve user"))
	}

	task, err := s.access.CreateTask(tsk)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, task)
}

func (s *Service) GetTaskByUserId(ctx echo.Context) error {

	from := ctx.QueryParam("from")
	to := ctx.QueryParam("to")

	f, err := strconv.Atoi(from)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	t, err := strconv.Atoi(to)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	userId, ok := ctx.Get("id").(int)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, errors.New("couldnt resolve user"))
	}

	result, err := s.access.GetTask(userId, f, t)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, result)
}

func hashPassword(raw string) (string, error) {

	dk, err := scrypt.Key([]byte(raw), []byte(viper.GetString("hashSecret")), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dk), nil
}
