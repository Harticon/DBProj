package DBproj

import (
	"fmt"
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

	fmt.Println(usr)

	//Valid email
	if !validateEmail(usr.Email) {
		fmt.Println("email not valid")
		return err
	}

	//Hash password

	//h := sha256.Sum256([]byte(usr.Password))
	//usr.Password = string(h[:])

	s.access.CreateUser(usr)
	return nil
}

func (s *Service) SignIn(ctx echo.Context) error {
	//Validate login
	//Validate passowrd
	//return login token

	s.access.GetUser()

	return nil
}

func (s *Service) SetTask(ctx echo.Context) error {

	//Get timestamp
	s.access.CreateTask()

	return nil
}

func (s *Service) GetTaskByUserId(ctx echo.Context) error {

	//Validate taskID
	s.access.GetTask()

	return nil
}
