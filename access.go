package DBproj

import (
	"encoding/base64"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"
)

func hashPassword(raw string) (string, error) {
	dk, err := scrypt.Key([]byte(raw), []byte("salt&pepper"), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dk), nil
}

type Access struct {
	db *gorm.DB
}

func NewAccess(db *gorm.DB) *Access {
	return &Access{
		db: db,
	}
}

type IAccesser interface {
	CreateUser(usr User)
	GetUser()
	CreateTask()
	GetTask()
}

func (a *Access) CreateUser(usr User) {

	usr.Password, _ = hashPassword(usr.Password)
	a.db.Create(&usr)

}

func (a *Access) GetUser() {
	//READ FROM DB
}

func (a *Access) CreateTask() {
	//WRITE to DB
}

func (a *Access) GetTask() {
	//READ from DB
}
