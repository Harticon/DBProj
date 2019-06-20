package DBproj

import "github.com/jinzhu/gorm"

type AccessMock struct {
	userDb map[int]User
	taskDb map[int]Task
}

func NewAccessMock() *AccessMock {
	return &AccessMock{
		userDb: map[int]User{
			0: {
				Model: gorm.Model{
					ID: 0,
				},
				Firstname: "vojta",
				Lastname:  "hromadka",
				Email:     "email@com.cz",
				Password:  "vojta",
			},
			1: {
				Model: gorm.Model{
					ID: 1,
				},
				Firstname: "v",
				Lastname:  "h",
				Email:     "e@com.cz",
				Password:  "pass",
			},
		},
		taskDb: map[int]Task{
			0: {
				Name:      "FirstTask",
				UserId:    0,
				ExecuteAt: 77,
			},
			1: {
				Name:      "SecondTask",
				UserId:    0,
				ExecuteAt: 108,
			},
			2: {
				Name:      "FirstTask",
				UserId:    1,
				ExecuteAt: 216,
			},
		},
	}
}

type IMockAccesser interface {
	CreateUser(usr User) (User, error)
	GetUser(email, password string) (User, error)
	CreateTask(tsk Task) (Task, error)
	GetTask(userId, f, t int) ([]Task, error)
}

func (a *AccessMock) CreateUser(usr User) (User, error) {

	a.userDb[2] = User{Email: usr.Email,
		Lastname:  usr.Lastname,
		Firstname: usr.Firstname,
		Password:  usr.Password}

	return usr, nil
}

func (a *AccessMock) GetUser(email, password string) (User, error) {
	for _, v := range a.userDb {
		if v.Email == email && v.Password == password {
			return v, nil
		}
	}
	return User{}, nil
}

func (a *AccessMock) CreateTask(tsk Task) (Task, error) {

	a.taskDb[3] = Task{
		ExecuteAt: tsk.ExecuteAt,
		UserId:    tsk.UserId,
		Name:      tsk.Name,
	}
	return tsk, nil
}

func (a *AccessMock) GetTask(userId, f, t int) ([]Task, error) {
	var tasks []Task

	for _, v := range a.taskDb {
		if userId == v.UserId && f <= v.ExecuteAt && t >= v.ExecuteAt {
			tasks = append(tasks, v)
		}
	}
	return tasks, nil

}
