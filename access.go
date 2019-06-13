package DBproj

import (
	"encoding/base64"
	"fmt"
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
	GetUser(usr User) User
	CreateTask(tsk Task)
	GetTask(userId int) []Task
}

func (a *Access) CreateUser(usr User) {

	usr.Password, _ = hashPassword(usr.Password)
	err := a.db.Create(&usr)
	if err.Error != nil {
		fmt.Println("email already registered")
	}

	// check email redundacy
}

func (a *Access) GetUser(usr User) User {
	usr.Password, _ = hashPassword(usr.Password)

	var query User

	err := a.db.Where("email = ? AND password = ?", usr.Email, usr.Password).Find(&query).Error

	if err != nil {
		fmt.Println("Invalid user/password!")
	}

	return query

	//READ FROM DB
}

func (a *Access) CreateTask(tsk Task) {

	err := a.db.Create(&tsk).Error
	if err != nil {
		fmt.Println("Couldnt write task to db")
	}

}

func (a *Access) GetTask(userId int) []Task {

	var query []Task
	err := a.db.Where("user_id = ?", userId).Find(&query).Error
	if err != nil {
		fmt.Println("You dont have any tasks yet")
	}

	return query
}
