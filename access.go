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
	CreateUser(usr User) error
	GetUser(usr User) User
	CreateTask(tsk Task)
	GetTask(userId, f, t int) []Task
}

func (a *Access) CreateUser(usr User) error {

	usr.Password, _ = hashPassword(usr.Password)
	err := a.db.Create(&usr)
	if err.Error != nil {
		fmt.Println(err.Error)
		return err.Error
	}

	// check email redundacy
	return nil
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

func (a *Access) GetTask(userId, f, t int) []Task {

	var query []Task
	fmt.Printf("%v %v", f, t)
	err := a.db.Where("user_id = ? AND execute_at <= ? AND execute_at >= ?", userId, t, f).Find(&query).Error

	if err != nil {
		fmt.Println("couldnt read from db")
		return []Task{}
	}
	return query
}
