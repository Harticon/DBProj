package DBproj

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type serviceSuite struct {
	suite.Suite
	service *Service

	echo *echo.Echo
	db   *gorm.DB
}

func (s *serviceSuite) SetupSuite() {

	viper.SetDefault("db.conn", "test.db")

	db, err := gorm.Open("sqlite3", viper.GetString("db.conn"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Task{})

	access := NewAccess(db)
	s.service = NewService(access)
	s.echo = echo.New()
	s.db = db

	user := &User{
		Firstname: "vojta",
		Lastname:  "hromadka",
		Email:     "hromadkavojta@gmail.com",
		Password:  "DfZmTUg5mGgzQATsDXlLv5c2gH9+kuNvjIaq9MyuxkU=",
	}

	err = s.db.Create(&user).Error
	s.Nil(err)

}

func (s *serviceSuite) SetupTest() {

}

func (s *serviceSuite) TearDownTest() {

}

func (s *serviceSuite) TearDownSuite() {

	s.db.DropTable(&User{})
	s.db.DropTable(&Task{})
	err := s.db.Close()
	s.Nil(err)

}

//---------------------------------------------------------------------------------------------------------------------
func (s *serviceSuite) getToken() string {
	user := &User{
		Email:    "hromadkavojta@gmail.com",
		Password: "vojta",
	}

	body, _ := json.Marshal(&user)
	reqs := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(string(body)))
	reqs.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recs := httptest.NewRecorder()
	c := s.echo.NewContext(reqs, recs)

	err := s.service.SignIn(c)
	s.Nil(err)

	var t Token

	err = json.Unmarshal(recs.Body.Bytes(), &t)
	s.Nil(err)

	return t.Token

}

//----------------------------------------------------------------------------------------------------------------------

func TestApiSuite(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}

func (s *serviceSuite) TestSignUp() {
	candidates := []struct {
		User         *User
		expectedCode int
		expectedErr  error
	}{
		{
			User: &User{
				Firstname: "vojta",
				Lastname:  "hromadka",
				Email:     "lolololololol@gmail.com",
				Password:  "vojta",
			},
			expectedCode: http.StatusCreated,
			expectedErr:  nil,
		},
		{
			User: &User{
				Firstname: "",
				Lastname:  "",
				Email:     "",
				Password:  "",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  nil,
		},
	}

	for i, candidate := range candidates {

		body, err := json.Marshal(&candidate.User)
		req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)

		err = s.service.SignUp(c)
		s.NoError(err)
		s.Equalf(candidate.expectedCode, rec.Code, "\n candidate: %d\n", i+1)
	}

}

func (s *serviceSuite) TestSignIn() {
	candidates := []struct {
		User         *User
		Token        string
		expectedCode int
		expectedErr  error
	}{
		{
			User: &User{
				Email:    "hromadkavojta@gmail.com",
				Password: "vojta",
			},
			expectedCode: http.StatusAccepted,
			expectedErr:  nil,
		},
		{
			User: &User{
				Email:    "non_existing_user",
				Password: "pass",
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  nil,
		},
	}

	for i, candidate := range candidates {

		body, err := json.Marshal(&candidate.User)
		s.NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := s.echo.NewContext(req, rec)

		err = s.service.SignIn(c)
		s.NoError(err)

		s.Equalf(candidate.expectedCode, rec.Code, "\n candidate: %d\n", i+1)

	}

}

func (s *serviceSuite) TestSetTask() {

	candidates := []struct {
		Task         *Task
		expectedCode int
		expectedErr  error
	}{
		{
			Task: &Task{
				Name:      "task1",
				ExecuteAt: 178,
			},
			expectedCode: http.StatusCreated,
			expectedErr:  nil,
		},
		{
			Task: &Task{
				Name:      "task32",
				ExecuteAt: 12,
			},
			expectedCode: http.StatusCreated,
			expectedErr:  nil,
		},
	}

	token := s.getToken()

	for i, candidate := range candidates {

		body, err := json.Marshal(&candidate.Task)
		req := httptest.NewRequest(http.MethodPost, "/task/create", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		ctx := s.echo.NewContext(req, rec)

		t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		s.Nil(err)

		claims, _ := t.Claims.(jwt.MapClaims)
		ctx.Set("id", int(claims["id"].(float64)))

		err = s.service.SetTask(ctx)
		s.Nil(err)
		s.Equalf(candidate.expectedCode, rec.Code, "\n candidate: %d\n", i+1)

	}

}

func (s *serviceSuite) TestGetTaskByUserId() {

	candidates := []struct {
		params       []string
		paramsVal    []string
		expectedCode int
		expectedErr  error
	}{
		{
			params:       []string{"from", "to"},
			paramsVal:    []string{"0", "250"},
			expectedCode: http.StatusOK,
			expectedErr:  nil,
		},
		{
			params:       []string{"from", "to"},
			paramsVal:    []string{"0aw", "250"},
			expectedCode: http.StatusBadRequest,
			expectedErr:  nil,
		},
		{
			params:       []string{"from", "to"},
			paramsVal:    []string{"0", "0"},
			expectedCode: http.StatusOK,
			expectedErr:  nil,
		},
		{
			params:       []string{"fromneco", "to"},
			paramsVal:    []string{"0", "250"},
			expectedCode: http.StatusBadRequest,
			expectedErr:  nil,
		},
	}

	token := s.getToken()

	for i, candidate := range candidates {

		f := make(url.Values)
		f.Set(candidate.params[0], candidate.paramsVal[0])
		f.Set(candidate.params[1], candidate.paramsVal[1])

		req := httptest.NewRequest(http.MethodGet, "/?"+f.Encode(), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		req.Header.Set(echo.HeaderAuthorization, token)
		rec := httptest.NewRecorder()
		ctx := s.echo.NewContext(req, rec)

		t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		s.Nil(err)

		claims, _ := t.Claims.(jwt.MapClaims)
		ctx.Set("id", int(claims["id"].(float64)))

		err = s.service.GetTaskByUserId(ctx)
		s.Nil(err)
		s.Equalf(candidate.expectedCode, rec.Code, "\n candidate: %d\n", i+1)

	}

}
