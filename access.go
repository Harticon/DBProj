package DBproj

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Access struct {
	db *gorm.DB
}

func NewAccess(db *gorm.DB) *Access {
	return &Access{
		db: db,
	}
}

type IAccesser interface {
	CreateUser(usr User) (User, error)
	GetUser(usr User) (User, error)
	CreateTask(tsk Task) (Task, error)
	GetTask(userId, f, t int) ([]Task, error)
}

func (a *Access) CreateUser(usr User) (User, error) {

	err := a.db.Create(&usr).Error
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	return usr, nil
}

func (a *Access) GetUser(usr User) (User, error) {
	var query User

	err := a.db.Where("email = ? AND password = ?", usr.Email, usr.Password).Find(&query).Error
	if err != nil {
		//fmt.Println("err:", err)
		return User{}, err
	}
	return query, nil
}

func (a *Access) CreateTask(tsk Task) (Task, error) {

	err := a.db.Create(&tsk).Error
	if err != nil {
		return Task{}, err
	}
	return tsk, nil

}

func (a *Access) GetTask(userId, f, t int) ([]Task, error) {
	var query []Task
	err := a.db.Where("user_id = ? AND execute_at <= ? AND execute_at >= ?", userId, t, f).Find(&query).Error
	if err != nil {
		fmt.Println("couldn't read from db")
		return []Task{}, err
	}
	return query, nil
}
