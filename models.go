package DBproj

type User struct {
	ID        int    `gorm:"AUTO_INCREMENT"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email" gorm:"unique_index;not null"`
	Password  string `json:"password" gorm:"not null"`
}

type Task struct {
	UserId    int
	Name      string `json:"name"`
	ExecuteAt string `json:"executeAt"`
}
