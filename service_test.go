package DBproj

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	userJSON = `{"firstname":"Jon","lastname":"Snow","email":"jon@labstack.com","password":"foo"}`
)

func TestNewAccess(t *testing.T) {

}

func TestNewService(t *testing.T) {

}

func TestUserMiddleware(t *testing.T) {

}

func TestService_GetTaskByUserId(t *testing.T) {

}

func TestService_SignIn(t *testing.T) {

}

func TestService_SetTask(t *testing.T) {

}

func TestService_SignUp(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Printf("%v", c)

}

func TestAccess_CreateUser(t *testing.T) {

}

func TestAccess_GetUser(t *testing.T) {

}

func TestAccess_CreateTask(t *testing.T) {

}

func TestAccess_GetTask(t *testing.T) {

}
